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

type Resource []string

//AllResources check whether this resource match all resources
func (r *Resource) matchAll() bool {
	if len(*r) == 1 && (*r)[0] == asterisk {
		return true
	}
	return false
}

// check whether resource math resource in context
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

func (r *Resource) evaluate(c *Context) error {
	for i, rs := range *r {
		if !containBrace(rs) {
			continue
		}
		// split resource string into tokens
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

// use context evaluate the value represented by token
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

// check whether some rule in rules match target
func match0(rules []string, target string) bool {
	tokens := tokenSplit.splitExpression(target)
	for _, rule := range rules {
		if match(tokenSplit.splitExpression(rule), tokens) {
			return true
		}
	}
	return false
}

// check whether rule match target
func match(rule Tokens, target Tokens) bool {

	// if rule is    a:b:c:*
	// and target is a:b
	// return false
	// authority requested in target is higher than authority in rule
	if len(target.split) < len(rule.split) {
		return false
	}
	for i, ps := range rule.split {
		// equal
		if (*rule.stringPointer)[ps[0]:ps[1]] == (*target.stringPointer)[target.split[i][0]:target.split[i][1]] {
			continue
		}
		// * match anything
		// /a/*/b match a/c/b , a/d/b ..
		if (*rule.stringPointer)[ps[0]:ps[1]] == "*" {
			continue
		}
		//[a b c] match a , b, c
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
	// authority requested in target is lower than authority in rule
	return true
}

func isListString(s string) bool {
	return len(s) > 2 && s[0] == openBrackets && s[len(s)-1] == closeBrackets
}
