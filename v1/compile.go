package v1

import (
	"fmt"
	"reflect"
	"strconv"
)

// compile use expression and context evaluate the result
// example:
// context:
//{
//	"user": {
//		"name": "tom",
//		"age": 1,
//		"mother": "tomMother"
//	},
//	"role": {
//		"name": "baby"
//	},
//	"resource": {
//		"name": "milk",
//		"owner": "tomMather"
//	},
//	"action": "drink"
//}
// {$.user.name}=tom
// {$.user.age}=1
// {$.resource.owner}="tomMather"

// a . b, a[b], a[b:c], a[*]
type op uint

const (
	value op = iota // a.b
	slice           // a[b:c]
	index           // a[b]
	scan            // a[*]
)

type step struct {
	op   op
	key  string        // context's value key or index of array
	args []interface{} // arguments of slice if [a:b] args is [a,b]
}

// object look path
type compiled struct {
	steps []*step
}

// lookup the value in context
func (c *compiled) lookup(context interface{}) (interface{}, error) {
	for _, s := range c.steps {
		var err error = nil
		// a.b
		if s.op == value {
			context, err = lookupByValue(context, s)
			if err != nil {
				return nil, err
			}
		}
		// a[b] a[b:c]
		if s.op == index {
			context, err = lookupByIndex(context, s)
			if err != nil {
				return nil, err
			}
		}
		// a[*]
		if s.op == scan {
			context, err = lookupByScan(context, s)
			if err != nil {
				return nil, err
			}
		}
	}
	return context, nil
}

// a:=[{k:1},{k:2}] --> a[*].k = [1,2]
func getScanValues(c interface{}, s *step) (interface{}, error) {
	var values []interface{}
	cv := reflect.ValueOf(c)
	for i := 0; i < cv.Len(); i++ {
		i2, err := lookupByValue(cv.Index(i).Interface(), s)
		if err != nil {
			return nil, err
		}
		values = append(values, i2)
	}
	return values, nil
}

// a=[{k:1},{k:2}], a[*]=[{k:1},{k:2}]
func lookupByScan(c interface{}, s *step) (interface{}, error) {
	cv := reflect.ValueOf(c)
	switch cv.Kind() {
	case reflect.Slice | reflect.Array:
		return cv.Interface(), nil
	}
	return nil, fmt.Errorf("can't get [*] from kind %s", cv.Kind())
}

func lookupByIndex(c interface{}, s *step) (interface{}, error) {
	cv := reflect.ValueOf(c)
	switch cv.Kind() {
	case reflect.Array:
		return cv.Index(s.args[0].(int)).Interface(), nil
	case reflect.Slice:
		if len(s.args) > 1 {
			return cv.Slice(s.args[0].(int), s.args[1].(int)).Interface(), nil
		}
		return cv.Index(s.args[0].(int)).Interface(), nil
	case reflect.Ptr:
		return lookupByIndex(c, s)
	}
	return nil, fmt.Errorf("can't get value %s from Kind %s", s.key, cv.Kind().String())
}

// return the value of `.` value operation
func lookupByValue(c interface{}, s *step) (interface{}, error) {
	cv := reflect.ValueOf(c)
	switch cv.Kind() {
	// if c is a map[string]interface{} return value of c[s.key]
	case reflect.Map:
		jsonMap := c.(map[string]interface{})
		return jsonMap[s.key], nil
	// if c is a struct use s.key as struct's field name or fields json tag
	case reflect.Struct:
		fileName, err := getStructFileByFiledNameOrJsonTag(s.key, reflect.TypeOf(c))
		if err != nil {
			return nil, err
		}
		return cv.FieldByName(fileName).Interface(), nil
	case reflect.Ptr:
		return lookupByValue(cv.Elem().Interface(), s)
	// if c is array or slice return every item in c computed by s
	case reflect.Array | reflect.Slice:
		return getScanValues(c, s)
	}

	return nil, fmt.Errorf("can't get value %s from Kind %s", s.key, cv.Kind().String())
}

func getStructFileByFiledNameOrJsonTag(n string, t reflect.Type) (string, error) {
	// get value by name
	if field, find := t.FieldByName(n); find {
		return field.Name, nil
	}
	// get value by json tag
	for i := 0; i < t.NumField(); i++ {
		if v, ok := t.Field(i).Tag.Lookup("json"); ok && v == n {
			return t.Field(i).Name, nil
		}
	}
	return "", fmt.Errorf("can't find %s in struct %+v", n, t)
}

func compile(exp string) (*compiled, error) {
	sequence, err := tokenize(exp)
	if err != nil {
		return nil, err
	}
	steps := make([]*step, 0, 10)
	for sequence.hasNext() {
		var s step
		token := sequence.pop()
		if token == "$" {
			continue
		}
		if token == "." {
			s = step{
				op:  value,
				key: sequence.pop(),
			}
		}
		if token == "[" {
			// [a
			p1 := sequence.pop()
			// [a: or [a]
			p2 := sequence.pop()
			if p2 == ":" {
				p1v, err := strconv.Atoi(p1)
				if err != nil {
					return nil, fmt.Errorf("except a number after [ got " + p1)
				}
				// [a:b]
				p3 := sequence.pop()
				p3v, err := strconv.Atoi(p3)
				if err != nil {
					return nil, fmt.Errorf("except a number after [: got " + p3)
				}
				s = step{
					op:   slice,
					args: []interface{}{p1v, p3v},
				}
			}
			if p2 == "]" {
				if p1 == "*" {
					s = step{
						op: scan,
					}
				} else {
					p1v, err := strconv.Atoi(p1)
					if err != nil {
						return nil, fmt.Errorf("except number or * after [ got " + p1)
					}
					s = step{
						op:   index,
						args: []interface{}{p1v},
					}
				}

			}
		}
		steps = append(steps, &s)
	}
	return &compiled{steps: steps}, nil
}
