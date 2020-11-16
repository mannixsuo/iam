package v1

import (
	"fmt"
	"testing"
)

func TestResource_Evaluate(t *testing.T) {
	ctx := Context{
		Requester: map[string]interface{}{"group": []string{"test", "test2"}, "name": "mmsuo"},
	}

	var r Resource = []string{"a:b:c/{$.requester.group}", "a:b:c/{$.requester.name}"}
	err := r.evaluate(&ctx)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r)
}

func BenchmarkResource_Evaluate(b *testing.B) {
	ctx := Context{
		Requester: map[string]interface{}{"group": "test"},
	}
	var r Resource = []string{"a:b:c/{$.user.group}"}
	for i := 0; i < b.N; i++ {
		r.evaluate(&ctx)
	}
}
