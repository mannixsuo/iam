package statement

import (
	"auth/auth"
	"auth/policy"
	"errors"
)

type Action struct {
	Action interface{} `json:"Action"`
}

func (a *Action) StatementType() policy.StatementType {
	return policy.ActionStatement
}

const allToken = "*"

//All 判断这个action是不是所有类型的action
func (a *Action) all() bool {
	if v, ok := a.Action.(string); ok {
		if v == allToken {
			return true
		}
	}
	return false
}

func (a *Action) actionList() (*[]string, error) {
	if a.all() {
		return nil, errors.New("this action represents all action")
	}
	if strings, ok := a.Action.([]interface{}); ok {
		v := make([]string, len(strings))
		for _, str := range strings {
			if s, ok := str.(string); ok {
				v = append(v, s)
			} else {
				return nil, errors.New("action value error,only string are allowed")
			}
		}
		return &v, nil
	}
	return nil, errors.New("format error Action = * | [action_string,]")
}

// 对比该action能否与context中的action匹配
func (a *Action) Match(c *auth.Context) (bool, error) {
	if a.all() {
		return true, nil
	}
	actionList, err := a.actionList()
	if err != nil {
		return false, err
	}
	for _, action := range *actionList {
		if action == c.Action {
			return true, nil
		}
	}
	return false, nil
}

func (a *Action) Evaluate(c *auth.Context) {

}
