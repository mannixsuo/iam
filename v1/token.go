package v1

import "fmt"

// like $ . a . b [ 1 ] . c
type tokenSequence struct {
	tokens Tokens //["$",".","a",".","b","[","1","]",".","c"]
	cIndex int    //current index
}

// split exp by defaultTokenSplits
// $.a.b[1].c => $ . a . b [ 1 ] . c
func tokenize(exp string) (*tokenSequence, error) {
	if exp[0] != '$' {
		return nil, fmt.Errorf("expression parser error: %s.expression should start with $", exp)
	}
	split := TokenSplit{Splits: defaultTokenSplits, SaveToken: true}

	sequence := tokenSequence{tokens: split.splitExpression(exp)}

	return &sequence, nil
}

type Tokens struct {
	stringPointer *string
	split         [][2]int
}

func (t *Tokens) append(bottom, head int) {
	if bottom == head {
		return
	}
	t.split = append(t.split, [2]int{bottom, head})
}

func (t Tokens) equals(other []string) bool {
	if len(t.split) == len(other) {
		for i := 0; i < len(other); i++ {
			if ((*(t.stringPointer))[t.split[i][0]:t.split[i][1]]) != other[i] {
				return false
			}
		}
		return true
	}
	return false
}

var defaultTokenSplits = []byte{byte('.'), byte('['), byte(']'), byte(':')}

// TokenSplit is used for split string in tokens by Splits
type TokenSplit struct {
	Splits    []byte
	SaveToken bool // save token when split string
}

// check whether split the char c
func (t *TokenSplit) shouldSplit(c byte) bool {
	for _, s := range t.Splits {
		if s == c {
			return true
		}
	}
	return false
}

// split exp by splits in tokenSplit
func (t *TokenSplit) splitExpression(exp string) Tokens {
	expLen := len(exp)
	var token = Tokens{stringPointer: &exp, split: make([][2]int, 0, expLen)}
	bottom := 0
	for head := 0; head < expLen; head++ {
		if t.shouldSplit(exp[head]) {
			// save token before
			token.append(bottom, head)
			// save split token
			if t.SaveToken {
				token.append(head, head+1)
			}
			// reset token start index
			bottom = head + 1
		}
	}
	// save the last token
	token.append(bottom, expLen)
	return token
}

// return the element on the top of tokenSequence
func (t *tokenSequence) pop() (token string) {
	token = (*(t.tokens.stringPointer))[t.tokens.split[t.cIndex][0]:t.tokens.split[t.cIndex][1]]
	if t.hasNext() {
		t.cIndex++
	}
	return
}

// back
func (t *tokenSequence) back() {
	t.cIndex--
}

func (t *tokenSequence) hasNext() bool {
	return t.cIndex < len(t.tokens.split)
}
