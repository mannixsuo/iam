package v1

// Context represent who want do what to which resources.
// such as requester want do action to resource with condition
type Context struct {
	Action    string                 `json:"action"`
	Requester map[string]interface{} `json:"requester"`
	Resource  string                 `json:"resource"`
	Condition map[string]interface{} `json:"condition"`
}
