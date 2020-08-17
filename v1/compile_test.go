package v1

import (
	"fmt"
	"testing"
)

func TestTokenSplitSplit(t *testing.T) {
	s := TokenSplit{SaveToken: false, Splits: []byte{byte('.'), byte('['), byte(']')}}
	tokens := s.split2tokens(".a.b[c].e")
	fmt.Println(tokens)
}

func TestLookup(t *testing.T) {
	context := struct {
		A string `json:"a"`
		B string `json:"b"`
	}{A: "aaa", B: "bbb"}
	c, _ := compile("$.b")
	lookup, err := c.lookup(context)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lookup)
	b := struct {
		A struct{ B string }
	}{A: struct{ B string }{B: "cccc"}}
	c, _ = compile("$.A.B")
	i, err := c.lookup(b)
	fmt.Println(i)
	d := make(map[string]interface{})
	d["A"] = "aaaaa"
	d["B"] = "bbbbb"
	c, _ = compile("$.B")
	lookup, _ = c.lookup(d)
	fmt.Println(lookup)
	user := make(map[string]interface{})
	user["name"] = "baby"
	user["age"] = 1
	user["roles"] = []map[string]interface{}{{"name": []string{"n1", "n2"}}, {"name": []string{"n3", "n4"}}, {"name": []string{"n5", "n6"}}}
	role1, _ := compile("$.roles[*].name")
	i2, err := role1.lookup(user)
	fmt.Println(fmt.Sprint(i2))
}

func TestGetIndex(t *testing.T) {
	a := make(map[string]interface{})
	a["a"] = []int{1, 2, 3}
	c, _ := compile("$.a[2]")
	lookup, _ := c.lookup(a)
	fmt.Println(lookup)
}

func Benchmark_compiled_lookup(b *testing.B) {
	a := make(map[string]interface{})
	a["a"] = []int{1, 2, 3}
	c, _ := compile("$.a[2]")
	for i := 0; i < b.N; i++ {
		_, _ = c.lookup(a)
	}
}
