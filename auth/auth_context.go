package auth

// Context 代表 who want do what to which resources.
type Context struct {
	Action    string                 `json:"action"`
	User      map[string]interface{} `json:"user"`
	Resource  string                 `json:"resource"`
	Condition map[string]interface{} `json:"condition"`
}
