package auth

// Context 代表 who want do what to which resources.
type Context struct {
	User     *User     `json:"user"`
	Resource *Resource `json:"resource"`
	Action   *Action   `json:"action"`
}

type User map[string]interface{}

type Resource string

type Action string
