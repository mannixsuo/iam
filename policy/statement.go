package policy

import "auth/auth"

// StatementType 声明类型 现在主要有4种
type StatementType string

const (
	ActionStatement    StatementType = "Action"
	EffectStatement    StatementType = "Effect"
	ResourceStatement  StatementType = "Resource"
	ConditionStatement StatementType = "Condition"
)

type Statement interface {
	//返回该声明的类型
	StatementType() StatementType
	//返回该声明是否和ctx中对应的属性匹配
	Match(c *auth.Context) (bool, error)
	//计算该声明的值
	Evaluate(c *auth.Context)
}
