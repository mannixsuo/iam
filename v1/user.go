package v1

// 用户 可以关联角色 关联组 关联策略
type User struct {
	Roles    []*Role   `json:"roles"`
	Groups   []*Group  `json:"groups"`
	Policies []*Policy `json:"policies"`
}

// 角色 可以关联策略
type Role struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Policies []*Policy `json:"policies"`
}

// 组 可以关联 策略
type Group struct {
	Policies []*Policy `json:"policies"`
}
