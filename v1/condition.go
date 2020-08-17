package v1

type Condition map[string][]interface{}


// Todo
func (c *Condition) Match(ac *Context) (bool, error) {
	return false, nil
}
