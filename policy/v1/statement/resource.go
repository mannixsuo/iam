package statement

import (
	"auth/policy"
	"errors"
)

type Resource struct {
	Resource interface{} `json:"Resource"`
}

func (r *Resource) StatementType() policy.StatementType {
	return policy.ResourceStatement
}

//AllResources 判断这个Resource是不是代表了所有资源
func (r *Resource) All() bool {
	if v, ok := r.Resource.(string); ok {
		if v == allToken {
			return true
		}
	}
	return false
}

func (r *Resource) ResourceList() (*[]string, error) {
	if r.All() {
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

