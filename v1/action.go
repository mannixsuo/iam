package v1

type Action []string

const allToken = "*"

//All 判断这个action是不是所有类型的action
func (a *Action) all() bool {
	if len(*a) == 1 && (*a)[0] == allToken {
		return true
	}
	return false
}

// 对比该action能否与context中的action匹配
func (a *Action) match(c *Context) (bool, error) {
	if a.all() {
		return true, nil
	}
	if match0(*a, c.Action) {
		return true, nil
	}
	return false, nil
}
