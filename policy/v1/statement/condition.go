package statement

import (
	"auth/auth"
	"auth/policy"
)

type Condition map[string][]interface{}

func (c *Condition) StatementType() policy.StatementType {
	return policy.ConditionStatement
}

// Todo
func (c *Condition) Match(ac *auth.Context) (bool, error) {
	return false, nil
}
