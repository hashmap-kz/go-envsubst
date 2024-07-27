package tok

import (
	"errors"
	"fmt"
	"github.com/hashmap.kz/go-envsubst/pkg/cbuf"
	"github.com/hashmap.kz/go-envsubst/pkg/util"
)

type TokType int

const (
	TokenTypeVar TokType = iota
	TokenTypeTxt
	TokenTypeEof
)

type Token struct {
	Value string
	Type  TokType
	Line  int
}

type Tokenlist struct {
	Tokens []*Token
}

func Tokenize(content string) (*Tokenlist, error) {
	b, err := cbuf.CBufNew(content)
	if err != nil {
		return nil, err
	}

	tokens := []*Token{}
	for {
		if b.IsEof() {
			break
		}
		t, _ := nex2(b)
		tokens = append(tokens, t)
	}

	return &Tokenlist{
		Tokens: tokens,
	}, nil
}

func nex2(b *cbuf.CBuf) (*Token, error) {
	peek3, _ := b.Peekc3()
	c1 := peek3[0]
	c2 := peek3[1]
	c3 := peek3[2]

	if c1 == '$' && c2 == '{' && util.IsIdentStart(c3) {
		b.Move(2)
		ident, err := getOneIdent(b)
		next, _ := b.Nextc()
		if next != '}' {
			return nil, errors.New("unclosed_brace")
		}
		return ident, err
	}

	if b.IsEof() {
		return &Token{
			Value: "",
			Type:  TokenTypeEof,
			Line:  b.Line,
		}, nil
	}

	n, err := b.Nextc()
	if err != nil {
		return nil, err
	}
	if n == cbuf.EOF {
		return &Token{
			Value: "",
			Type:  TokenTypeEof,
			Line:  b.Line,
		}, nil
	}

	return &Token{
		Value: string(n),
		Type:  TokenTypeTxt,
	}, nil
}

func getOneIdent(b *cbuf.CBuf) (*Token, error) {
	value := ""
	line := b.Line
	for {
		peek, err := b.Peekc1()
		if err != nil {
			return nil, err
		}
		if !util.IsIdentTail(peek) {
			break
		}

		next := b.Next()
		value += next
	}
	return &Token{
		Value: value,
		Type:  TokenTypeVar,
		Line:  line,
	}, nil
}

func (tl *Tokenlist) DumpStat() {
	fmt.Println("LINE    NAME")
	for _, t := range tl.Tokens {
		if t == nil {
			break
		}
		if t.Type == TokenTypeEof {
			break
		}
		if t.Type == TokenTypeVar {
			fmt.Printf("%-7d %s\n", t.Line, t.Value)
		}
	}
}

func (tl *Tokenlist) Dump() {
	for _, t := range tl.Tokens {
		if t == nil {
			break
		}
		if t.Type == TokenTypeEof {
			break
		}
		if t.Type == TokenTypeVar {
			fmt.Printf("${%s}", t.Value)
		} else {
			fmt.Printf("%s", t.Value)
		}
	}
}
