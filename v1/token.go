package v1

import "fmt"

// tokenSequence 代表分割后的token队列
// 比如 $ . a . b [ 1 ] . c
type tokenSequence struct {
	tokens Tokens //["$",".","a",".","b","[","1","]",".","c"]
	cIndex int    //当前指针所在位置
}

// 使用defaultTokenSplits 将表达式exp进行分割
// $ .a .b [1] . c
func tokenize(exp string) (*tokenSequence, error) {
	if exp[0] != '$' {
		return nil, fmt.Errorf("expression parser error: %s.expression should start with $", exp)
	}
	split := TokenSplit{Splits: defaultTokenSplits, SaveToken: true}

	sequence := tokenSequence{tokens: split.splitExpression(exp)}

	return &sequence, nil
}

type Tokens []string

func (t *Tokens) append(s string) {
	if s != "" {
		*t = append(*t, s)
	}
}

func (t Tokens) equals(other []string) bool {
	if len(t) == len(other) {
		for i := 0; i < len(other); i++ {
			if t[i] != other[i] {
				return false
			}
		}
		return true
	}
	return false
}

var defaultTokenSplits = []byte{byte('.'), byte('['), byte(']'), byte(':')}

// 根据 Splits 来分割字符串
type TokenSplit struct {
	Splits    []byte
	SaveToken bool //是否将分隔符也保存
}

// 判断是否分割
func (t *TokenSplit) shouldSplit(c byte) bool {
	for _, s := range t.Splits {
		if s == c {
			return true
		}
	}
	return false
}

// 根据 TokenSplit 的splits字符分割字符串
func (t *TokenSplit) splitExpression(exp string) Tokens {
	var token = Tokens{}
	expLen := len(exp)
	bottom := 0
	for head := 0; head < expLen; head++ {
		if t.shouldSplit(exp[head]) {
			// 保存分隔符前的token
			token.append(exp[bottom:head])
			// 保存分隔符
			if t.SaveToken {
				token.append(exp[head : head+1])
			}
			// 重置token起始位置
			bottom = head + 1
		}
	}
	// 保存最后一个token
	token.append(exp[bottom:expLen])
	return token
}

// 返回队列最前面的元素
func (t *tokenSequence) pop() (token string) {
	token = t.tokens[t.cIndex]
	if t.hasNext() {
		t.cIndex++
	}
	return
}

// 回退一步
func (t *tokenSequence) back() {
	t.cIndex--
}

// 是否还剩token未读取
func (t *tokenSequence) hasNext() bool {
	return t.cIndex < len(t.tokens)
}
