package model

// action <resource>:<action>
type Action struct {
	Action []string
}

// system:<service-name>:<region>:<tenant-id>:<relative-id>
//
// support policy variable
// ${username}
type Resource struct {
	Resource []string
}

// effect allow | deny
type Statement struct {
	Id       int      `json:"-"`
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
}

//
type Policy struct {
	Id        int          `json:"-"`
	Version   string       `json:"Version"`
	Statement []*Statement `json:"Statements"`
}

type Group struct {
	Id   int64
	Name string
}

type Role struct {
	Id   int64
	Name string
}
