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
func (a *Action) All() bool {
	if v, ok := a.Action.(string); ok {
		if v == allToken {
			return true
		}
	}
	return false
}

func (a *Action) ActionList() (*[]string, error) {
	if a.All() {
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

func (a *Action) MatchContext(c *auth.Context) (bool, error) {
	if a.All() {
		return true, nil
	}
	actionList, err := a.ActionList()
	if err != nil {
		return false, err
	}
	for _, action := range *actionList {

	}
}
