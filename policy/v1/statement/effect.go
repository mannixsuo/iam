package statement

import (
	"auth/auth"
	"auth/policy"
)

// Effect 表示这个策略代表的是允许还是拒绝
// Allow 或者 Deny
type Effect string

const (
	Allow Effect = "Allow"
	Deny  Effect = "Deny"
)

func (e *Effect) StatementType() policy.StatementType {
	return policy.EffectStatement
}

func (e *Effect) Match(c *auth.Context) (bool, error) {
	return string(*e) == c.Action, nil
}
func (e *Effect) Evaluate(c *auth.Context){}