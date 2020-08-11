package v1

import (
	"reflect"
)

// compile 根据表达式和context 来计算结果
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
func compile(exp string, c interface{}) (string, error) {

	switch reflect.ValueOf(c).Kind() {
	case reflect.Map:

	}
}

// a . b, a[b]
type op uint

const (
	value op = iota //某个对象的属性
	index           //list中的某个对象
)

type token struct {
	op  op
	key interface{} // 属性的名称或者index
}

func getTokens(exp string) []*token {
	for i, t := range exp {
		if t == '.' {

		}
		if t == '[' {
		}
	}
}

func parseToken(exp string) {}
