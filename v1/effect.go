package v1

// Effect represents the policy allow or deny do action on resources
type Effect string

const (
	Allow Effect = "Allow"
	Deny  Effect = "Deny"
)
