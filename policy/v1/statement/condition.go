package statement

import "auth/policy"

type Condition map[string][]interface{}

func (c *Condition) StatementType() policy.StatementType {
	return policy.ConditionStatement
}
