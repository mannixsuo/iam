package v1

import (
	"fmt"
	"testing"
)

func TestTokens_equals(t *testing.T) {
	s := ""
	tokens := Tokens{stringPointer: &s}
	if !tokens.equals([]string{}) {
		t.Error()
	}
	s = "a"
	if !tokens.equals([]string{"a"}) {
		t.Error()
	}
	s = "ab"
	if !tokens.equals([]string{"a", "b"}) {
		t.Error()
	}
}

func TestTokenSplitSplit(t *testing.T) {
	s := TokenSplit{SaveToken: true, Splits: []byte{byte('.'), byte('['), byte(']')}}
	exp := ".a.b[c].e"

	tokens := s.splitExpression(exp)
	if !tokens.equals([]string{".", "a", ".", "b", "[", "c", "]", ".", "e"}) {
		t.Error()
	}
	s.SaveToken = false
	tokens = s.splitExpression(exp)
	if !tokens.equals([]string{"a", "b", "c", "e"}) {
		t.Error()
	}

	exp = ""
	tokens = s.splitExpression(exp)
	if !tokens.equals([]string{}) {
		t.Error()
	}

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
	role1, _ := compile("$.roles[*].name[2][1]")
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

func Test_lookupByValue(t *testing.T) {
	// map
	m := map[string]interface{}{
		"name": "tom",
		"age":  12,
	}
	if byValue, err := lookupByValue(m, &step{op: value, key: "name"}); err != nil {
		t.Error()
	} else if byValue.(string) != m["name"] {
		t.Error()
	}
	if byValue, err := lookupByValue(m, &step{op: value, key: "age"}); err != nil {
		t.Error()
	} else if byValue.(int) != m["age"] {
		t.Error()
	}
	// struct
	st := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{Name: "tom", Age: 10}

	if byValue, err := lookupByValue(st, &step{op: value, key: "name"}); err != nil {
		t.Error()
	} else if byValue.(string) != st.Name {
		t.Error()
	}
	if byValue, err := lookupByValue(st, &step{op: value, key: "age"}); err != nil {
		t.Error()
	} else if byValue.(int) != st.Age {
		t.Error()
	}
	// pointer
	if byValue, err := lookupByValue(&st, &step{op: value, key: "name"}); err != nil {
		t.Error()
	} else if byValue.(string) != st.Name {
		t.Error()
	}
	if byValue, err := lookupByValue(&st, &step{op: value, key: "age"}); err != nil {
		t.Error()
	} else if byValue.(int) != st.Age {
		t.Error()
	}
	// array slice
	a := []interface{}{st, m}
	if byValue, err := lookupByValue(a, &step{op: value, key: "age"}); err != nil {
		t.Error()
	} else {
		ints := byValue.([]interface{})
		if ints[0].(int) == st.Age && ints[1].(int) == m["age"].(int) {
		} else {
			t.Error()
		}
	}

	if byValue, err := lookupByValue(a[1:], &step{op: value, key: "age"}); err != nil {
		t.Error()
	} else {
		ints := byValue.([]interface{})
		if ints[0].(int) == m["age"].(int) {
		} else {
			t.Error()
		}
	}
}
