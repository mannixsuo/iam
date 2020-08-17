package v1

import (
	"fmt"
	"testing"
)

func TestResource_Evaluate(t *testing.T) {
	ctx := Context{
		User: map[string]interface{}{"group": "test", "name": "mmsuo"},
	}

	var r Resource = []string{"a:b:c/{$.user.group}", "a:b:c/{$.user.name}"}
	r.evaluate(&ctx)

	fmt.Println(r)
}

func BenchmarkResource_Evaluate(b *testing.B) {
	ctx := Context{
		User: map[string]interface{}{"group": "test"},
	}
	var r Resource = []string{"a:b:c/{$.user.group}"}
	for i := 0; i < b.N; i++ {
		r.evaluate(&ctx)
	}
}
