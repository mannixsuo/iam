package v1

import (
	"fmt"
	"strings"
)

const (
	openBraces    = "{"
	closeBraces   = "}"
	openBrackets  = '['
	closeBrackets = ']'
	colon         = ":"
	period        = "."
	asterisk      = "*"
)

// * or  ["acs:ecs:*:*:instance/inst-001", "acs:ecs:*:*:instance/inst-002", "acs:oss:*:*:mybucket", "acs:oss:*:*:mybucket/*"]
// string or []string

type Resource []string

//AllResources 判断这个Resource是不是代表了所有资源
func (r *Resource) matchAll() bool {
	if len(*r) == 1 && (*r)[0] == asterisk {
		return true
	}
	return false
}

// 资源是否和context中的匹配
func (r *Resource) match(c *Context) (bool, error) {
	if r.matchAll() {
		return true, nil
	}
	err := r.evaluate(c)
	if err != nil {
		return false, err
	}
	if match0(*r, c.Resource) {
		return true, nil
	}
	return false, nil
}

var tokenSplit = TokenSplit{
	Splits:    []byte{':', '/'},
	SaveToken: false,
}

var saveTokenSp = TokenSplit{
	Splits:    []byte{':', '/'},
	SaveToken: true,
}

type CompiledResource struct {
	Compiled []Tokens
}

// 将 resource 中带有 {} 的计算出来
func (r *Resource) evaluate(c *Context) error {
	for i, rs := range *r {
		if !containBrace(rs) {
			continue
		}
		// 分割资源字符串
		tokens := saveTokenSp.splitExpression(rs)
		for _, split := range tokens.split {
			if isExpression((*tokens.stringPointer)[split[0]:split[1]]) {
				s, err := evaluate((*tokens.stringPointer)[split[0]:split[1]], c)
				if err != nil {
					return err
				}
				(*r)[i] = strings.ReplaceAll((*r)[i], (*tokens.stringPointer)[split[0]:split[1]], fmt.Sprint(s))
			}
		}
	}
	return nil
}

func containBrace(resource string) bool {
	return strings.Contains(resource, openBraces)
}

func isExpression(exp string) bool {
	if len(exp) > 0 && exp[0] == '{' && exp[len(exp)-1] == '}' {
		return true
	}
	return false
}

// 计算token的实际值 返回是查询到的对象在fmt.sprint中的格式
func evaluate(token string, c *Context) (interface{}, error) {
	//{$.a.b.c}
	compile, err := compile(token[1 : len(token)-1])
	if err != nil {
		return "", err
	}
	lookup, err := compile.lookup(c)
	if err != nil {
		return "", err
	}
	return lookup, nil

}

// 判断rules中的规则是否匹配 target
func match0(rules []string, target string) bool {
	tokens := tokenSplit.splitExpression(target)
	for _, rule := range rules {
		if match(tokenSplit.splitExpression(rule), tokens) {
			return true
		}
	}
	return false
}

// 判断规则rule 是否匹配规则 target
func match(rule Tokens, target Tokens) bool {

	// rule    a:b:c:*
	// target  a:b
	// false
	// 说明target的权限 比rule中的权限要高
	if len(target.split) < len(rule.split) {
		return false
	}
	for i, ps := range rule.split {
		if (*rule.stringPointer)[ps[0]:ps[1]] == (*target.stringPointer)[target.split[i][0]:target.split[i][1]] {
			continue
		}
		// * 匹配任何值
		if (*rule.stringPointer)[ps[0]:ps[1]] == "*" {
			continue
		}
		//[a b c]
		if isListString((*rule.stringPointer)[ps[0]:ps[1]]) {
			if strings.Contains((*rule.stringPointer)[ps[0]:ps[1]], (*target.stringPointer)[ps[0]:ps[1]]) {
				continue
			}
		}
		return false
	}
	// rule a:b:c
	// target a:b:c or a:b:c:d
	// true
	// 权限匹配或者是target权限比rule低
	return true
}

func isListString(s string) bool {
	return len(s) > 2 && s[0] == openBrackets && s[len(s)-1] == closeBrackets
}
