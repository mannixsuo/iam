package statement

import "auth/policy"

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
