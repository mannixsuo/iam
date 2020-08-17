package v1

import (
	"fmt"
	"testing"
)

func TestPolicy_Allow(t *testing.T) {

	p := Policy{
		Version: 1,
		Statements: []*Statement{
			{
				Action:   &Action{"food:eat"},
				Resource: &Resource{"{$.user.name}:food:*"},
				Effect:   Allow,
			},
			{
				Action:   &Action{"toy:eat"},
				Resource: &Resource{"{$.user.name}:toy:*"},
				Effect:   Deny,
			},
		},
	}
	ctx := &Context{
		Action:   "food:*",
		User:     map[string]interface{}{"name": "tom"},
		Resource: "tom:food:bread",
	}
	ctx2 := &Context{
		Action:   "toy:eat",
		User:     map[string]interface{}{"name": "tom"},
		Resource: "tom:toy:car",
	}

	allow, m, err := p.Allow(ctx)
	fmt.Println(allow, m, err)
	a, b, err := p.Allow(ctx2)
	fmt.Println(a, b, err)
}

func BenchmarkPolicy_Allow(b *testing.B) {
	b.ReportAllocs()
	p := Policy{
		Version: 1,
		Statements: []*Statement{
			{
				Action:   &Action{"food:eat"},
				Resource: &Resource{"{$.user.name}:food:*"},
				Effect:   Allow,
			},
			{
				Action:   &Action{"toy:eat"},
				Resource: &Resource{"{$.user.name}:toy:*"},
				Effect:   Deny,
			},
		},
	}
	ctx := &Context{
		Action:   "food:*",
		User:     map[string]interface{}{"name": "tom"},
		Resource: "tom:food:bread",
	}
	for i := 0; i < b.N; i++ {
		p.Allow(ctx)
	}
}
