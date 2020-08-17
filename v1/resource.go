package v1

import (
	"fmt"
	"strings"
)

// * or  ["acs:ecs:*:*:instance/inst-001", "acs:ecs:*:*:instance/inst-002", "acs:oss:*:*:mybucket", "acs:oss:*:*:mybucket/*"]
// string or []string

type Resource []string

//AllResources 判断这个Resource是不是代表了所有资源
func (r *Resource) all() bool {
	if len(*r) == 1 && (*r)[0] == allToken {
		return true
	}
	return false
}

// 资源是否和context中的匹配
func (r *Resource) match(c *Context) (bool, error) {
	if r.all() {
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

var sp = TokenSplit{
	Splits:    []byte{':', '/'},
	SaveToken: false,
}

var saveTokenSp = TokenSplit{
	Splits:    []byte{':', '/'},
	SaveToken: true,
}

// 将 resource 中带有 {} 的计算出来
func (r *Resource) evaluate(c *Context) error {
	for i, rs := range *r {
		tokens := saveTokenSp.split2tokens(rs)
		for j, t := range tokens {
			s, err := evaluate(t, c)
			if err != nil {
				return err
			}
			tokens[j] = s
		}
		(*r)[i] = join(tokens)
	}
	return nil
}

func join(s []string) string {
	sb := strings.Builder{}
	for _, ss := range s {
		sb.WriteString(ss)
	}
	return sb.String()
}

// 计算token的实际值 返回是查询到的对象在fmt.sprint中的格式
func evaluate(token string, c *Context) (string, error) {
	tl := len(token)
	//{$.a.b.c}
	if token[0] == '{' && token[tl-1] == '}' {
		compile, err := compile(token[1 : tl-1])
		if err != nil {
			return "", err
		}
		lookup, err := compile.lookup(c)
		if err != nil {
			return "", err
		}
		return fmt.Sprint(lookup), nil
	}
	return token, nil
}

func match0(p []string, r string) bool {
	rs := sp.split2tokens(r)
	for _, ps := range p {
		if match(sp.split2tokens(ps), rs) {
			return true
		}
	}
	return false
}

func match(p []string, r []string) bool {
	pl := len(p) - 1
	rl := len(r) - 1
	// p  a:b:c:*
	// r  a:b
	// false
	if rl < pl {
		return false
	}
	for i, ps := range p {
		if ps == "*" {
			continue
		}
		if ps == r[i] {
			continue
		}
		return false
	}
	// p a:b:c
	// r a:b:c or a:b:c:d
	// true
	return true
}
