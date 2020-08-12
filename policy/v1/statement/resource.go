package statement

import (
	"auth/auth"
	"auth/policy"
	v1 "auth/policy/v1"
	"errors"
)

type Resource struct {
	// * or  ["acs:ecs:*:*:instance/inst-001", "acs:ecs:*:*:instance/inst-002", "acs:oss:*:*:mybucket", "acs:oss:*:*:mybucket/*"]
	// string or *[]string
	Resource interface{} `json:"Resource"`
}

func (r *Resource) StatementType() policy.StatementType {
	return policy.ResourceStatement
}

//AllResources 判断这个Resource是不是代表了所有资源
func (r *Resource) all() bool {
	if v, ok := r.Resource.(string); ok {
		if v == allToken {
			return true
		}
	}
	return false
}

func (r *Resource) resourceList() (*[]string, error) {
	if r.all() {
		return nil, errors.New("this resource represents all resources")
	}
	if resources, ok := r.Resource.([]interface{}); ok {
		v := make([]string, 0)
		for _, str := range resources {
			if s, ok := str.(string); ok {
				v = append(v, s)
			} else {
				return nil, errors.New("resource value error,only string are allowed")
			}
		}
		return &v, nil
	}
	return nil, errors.New("format error Resource = * | [resource,resource..]")
}

// 资源是否和context中的匹配
func (r *Resource) Match(c *auth.Context) (bool, error) {
	if r.all() {
		return true, nil
	}
	list, err := r.resourceList()
	if err != nil {
		return false, err
	}
	for _, res := range *list {
		if c.Resource == res {
			return true, nil
		}
	}
	return false, nil
}

func (r *Resource) Evaluate(c *auth.Context) {
	if !r.all() {
		sp := v1.TokenSplit{
			Splits:    []byte{':', '/'},
			SaveToken: false,
		}
		res := r.Resource.(*[]string)
		// [a:b:c:*:d:{$.e.f}:g/h/i/*,a:b:c:*:d:{$.e.f}:g/h/i/*]
		for _, exp := range *res {
			_ = sp.Split2tokens(exp)//todo

		}
	}
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
