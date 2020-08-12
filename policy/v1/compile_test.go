package v1

import (
	"fmt"
	"testing"
)

func TestTokenSplitSplit(t *testing.T) {
	s := TokenSplit{[]byte{byte('.'), byte('['), byte(']')}}
	tokens := s.Split2tokens(".a.b[c].e")
	fmt.Println(tokens)
}

func TestLookup(t *testing.T) {
	context := struct {
		A string `json:"a"`
		B string `json:"b"`
	}{A: "aaa", B: "bbb"}
	c, _ := Compile("$.b")
	lookup, err := c.Lookup(context)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lookup)
	b := struct {
		A struct{ B string }
	}{A: struct{ B string }{B: "cccc"}}
	c, _ = Compile("$.A.B")
	i, err := c.Lookup(b)
	fmt.Println(i)
	d := make(map[string]interface{})
	d["A"] = "aaaaa"
	d["B"] = "bbbbb"
	c, _ = Compile("$.B")
	lookup, _ = c.Lookup(d)
	fmt.Println(lookup)
	user := make(map[string]interface{})
	user["name"] = "baby"
	user["age"] = 1
	user["roles"] = []map[string]interface{}{{"name": "role1"}, {"name": "role2"}, {"name": []string{"n1", "n2"}}}
	role1, _ := Compile("$.roles[2].name")
	i2, err := role1.Lookup(user)
	fmt.Println(fmt.Sprint(i2))
}

func TestGetIndex(t *testing.T) {
	a := make(map[string]interface{})
	a["a"] = []int{1, 2, 3}
	c, _ := Compile("$.a[2]")
	lookup, _ := c.Lookup(a)
	fmt.Println(lookup)
}

func Benchmark_compiled_lookup(b *testing.B) {
	a := make(map[string]interface{})
	a["a"] = []int{1, 2, 3}
	c, _ := Compile("$.a[2]")
	for i := 0; i < b.N; i++ {
		_, _ = c.Lookup(a)
	}
}
