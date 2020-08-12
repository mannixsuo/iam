package v1

import (
	"fmt"
	"reflect"
	"strconv"
)

// Compile 根据表达式和context 来计算结果
// example:
// context:
//{
//	"user": {
//		"name": "tom",
//		"age": 1,
//		"mother": "tomMather"
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

// a . b, a[b], a[b:c]
type op uint

const (
	value op = iota //a.b 某个对象的属性
	split           //a[b:c] 切片
	index           //a[b] list中的某个对象
	scan            //todo a[*] 所有对象
)

type step struct {
	op   op
	key  string // 属性的名称或者index
	args []interface{}
}

type Compiled struct {
	steps []*step
}

func (c *Compiled) Lookup(context interface{}) (interface{}, error) {
	for _, s := range c.steps {
		var err error = nil
		if s.op == value {
			context, err = getValue(context, s)
			if err != nil {
				return nil, err
			}
		}
		if s.op == index {
			context, err = getIndex(context, s)
			if err != nil {
				return nil, err
			}
		}
	}
	return context, nil
}

func getIndex(c interface{}, s *step) (interface{}, error) {
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
		return getIndex(c, s)
	}
	return nil, fmt.Errorf("can't get value %s from Kind %s", s.key, cv.Kind().String())
}

//如果c是map,s.key作为map的key来获取值
//如果c是结构体,s.key作为结构体属性名称或者属性的json tag定义的名称来取值
func getValue(c interface{}, s *step) (interface{}, error) {
	cv := reflect.ValueOf(c)
	switch cv.Kind() {
	case reflect.Map:
		jsonMap := c.(map[string]interface{})
		return jsonMap[s.key], nil
	case reflect.Struct:
		tag, err := getStructFileByFiledNameOrJsonTag(s.key, reflect.TypeOf(c))
		if err != nil {
			return nil, err
		}
		return cv.FieldByName(tag).Interface(), nil
	case reflect.Ptr:
		return getValue(cv.Elem().Interface(), s)
	}
	return nil, fmt.Errorf("can't get value %s from Kind %s", s.key, cv.Kind().String())
}

func getStructFileByFiledNameOrJsonTag(n string, t reflect.Type) (string, error) {
	// 首先根据名称查找
	if field, find := t.FieldByName(n); find {
		return field.Name, nil
	}
	// 根据json tag查找
	for i := 0; i < t.NumField(); i++ {
		if v, ok := t.Field(i).Tag.Lookup("json"); ok && v == n {
			return t.Field(i).Name, nil
		}
	}
	return "", fmt.Errorf("can't find %s in struct %+v", n, t)
}

func Compile(exp string) (*Compiled, error) {
	sequence, err := tokenize(exp)
	if err != nil {
		return nil, err
	}
	steps := make([]*step, 0)
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
			p1v, err := strconv.Atoi(p1)
			if err != nil {
				return nil, fmt.Errorf("except number after [ got " + p1)
			}
			// [a: or [a]
			p2 := sequence.pop()
			if p2 == ":" {
				// [a:b]
				p3 := sequence.pop()
				p3v, err := strconv.Atoi(p3)
				if err != nil {
					return nil, fmt.Errorf("except number after [: got " + p3)
				}
				s = step{
					op:   split,
					args: []interface{}{p1v, p3v},
				}
			}
			if p2 == "]" {
				s = step{
					op:   index,
					args: []interface{}{p1v},
				}
			}
		}
		steps = append(steps, &s)
	}
	return &Compiled{steps: steps}, nil
}

// $ .a .b [1] . c
func tokenize(exp string) (*tokenSequence, error) {
	if exp[0] != '$' {
		return nil, fmt.Errorf("expression parser error: %s.expression should start with $", exp)
	}
	split := TokenSplit{Splits: defaultTokenSplits, SaveToken: true}

	sequence := tokenSequence{tokens: split.Split2tokens(exp)}

	return &sequence, nil
}

// $ . a . b [ 1 ] . c
type tokenSequence struct {
	tokens Tokens //["$",".","a",".","b","[","1","]",".","c"]
	cIndex int    //当前指针所在位置
}

type Tokens []string

func (t *Tokens) append(s string) {
	if s != "" {
		*t = append(*t, s)
	}
}

var defaultTokenSplits = []byte{byte('.'), byte('['), byte(']'), byte(':')}

// 根据 Splits 来分割字符串
type TokenSplit struct {
	Splits    []byte
	SaveToken bool //是否将分隔符也保存
}

// 判断是否分割
func (t *TokenSplit) shouldSplit(c byte) bool {
	for _, s := range t.Splits {
		if s == c {
			return true
		}
	}
	return false
}

func (t *TokenSplit) Split2tokens(exp string) Tokens {
	var currentToken []byte
	var token = Tokens{}
	for _, b := range []byte(exp) {
		if t.shouldSplit(b) {
			// 保存分隔符前的token,以及分割符
			token.append(string(currentToken))
			if t.SaveToken {
				token.append(string(b))
			}
			currentToken = make([]byte, 0)
			continue
		}
		currentToken = append(currentToken, b)
	}
	// 保存最后一个token
	token.append(string(currentToken))
	return token
}

// 返回队列最前面的元素
func (t *tokenSequence) pop() (token string) {
	token = t.tokens[t.cIndex]
	if t.hasNext() {
		t.cIndex++
	}
	return
}

// 回退一步
func (t *tokenSequence) back() {
	t.cIndex--
}

// 是否还剩token未读取
func (t *tokenSequence) hasNext() bool {
	return t.cIndex < len(t.tokens)
}

func parseToken(exp string) {}
