package v1

type Action []string

const allToken = "*"

//All check whether this action match all action
func (a *Action) matchAll() bool {
	if len(*a) == 1 && (*a)[0] == allToken {
		return true
	}
	return false
}

// check whether this actions match actions in context
func (a *Action) match(c *Context) (bool, error) {
	if a.matchAll() {
		return true, nil
	}
	if match0(*a, c.Action) {
		return true, nil
	}
	return false, nil
}
