package v1

type Statement struct {
	Action   *Action   `json:"Action"`
	Resource *Resource `json:"Resource"`
	Effect   Effect   `json:"Effect"`
}

func (p *Statement) match(c *Context) (bool, error) {
	am, err := p.Action.match(c)
	if err != nil {
		return false, err
	}
	pm, err := p.Resource.match(c)
	if err != nil {
		return false, err
	}
	if am && pm {
		return true, nil
	}
	return false, nil
}
