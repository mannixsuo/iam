package policy

// StatementType 声明类型 现在主要有4种
type StatementType string

const (
	ActionStatement    StatementType = "Action"
	EffectStatement    StatementType = "Effect"
	ResourceStatement  StatementType = "Resource"
	ConditionStatement StatementType = "Condition"
)

type Statement interface {
	StatementType() StatementType
}
