package tok

import (
	"errors"
	"fmt"
	"github.com/hashmap.kz/go-envsubst/pkg/cbuf"
	"github.com/hashmap.kz/go-envsubst/pkg/util"
	"os"
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

	var prevToken *Token
	for {
		if b.IsEof() {
			break
		}
		t, _ := nex2(b)

		// merge uninteresting for us 'plain text'
		if t.Type == TokenTypeTxt && (prevToken != nil && prevToken.Type == TokenTypeTxt) {
			if len(tokens) > 0 {
				// insert
				lastIdx := len(tokens) - 1
				tokens[lastIdx].Value += t.Value
			} else {
				// append
				tokens = append(tokens, t)
			}
		} else {
			tokens = append(tokens, t)
		}
		prevToken = t

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

	// ${EDITOR}
	if c1 == '$' && c2 == '{' && util.IsIdentStart(c3) {
		b.Move(2)
		ident, err := getOneIdent(b)
		next, _ := b.Nextc()
		if next != '}' {
			return nil, errors.New("unclosed_brace")
		}
		return ident, err
	}

	// $EDITOR
	if c1 == '$' && util.IsIdentStart(c2) {
		b.Move(1)
		ident, err := getOneIdent(b)
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
	fmt.Println("LINE    NAME                             VALUE")
	for _, t := range tl.Tokens {
		if t == nil {
			break
		}
		if t.Type == TokenTypeEof {
			break
		}
		if t.Type == TokenTypeVar {
			fmt.Printf("%-7d %-32s %s\n", t.Line, t.Value, os.Getenv(t.Value))
		}
	}
}

func (tl *Tokenlist) DumpRawUnexpanded() string {
	result := ""
	for _, t := range tl.Tokens {
		if t == nil {
			break
		}
		if t.Type == TokenTypeEof {
			break
		}
		if t.Type == TokenTypeVar {
			result += fmt.Sprintf("${%s}", t.Value)
		} else {
			result += fmt.Sprintf("%s", t.Value)
		}
	}
	return result
}

func (tl *Tokenlist) DumpExpanded() string {
	result := ""
	for _, t := range tl.Tokens {
		if t == nil {
			break
		}
		if t.Type == TokenTypeEof {
			break
		}
		if t.Type == TokenTypeVar {
			result += expandOneVar(t)
		} else {
			result += fmt.Sprintf("%s", t.Value)
		}
	}
	return result
}

// TODO: this may be optimized a LOT with the global hashtable

func expandOneVar(t *Token) string {
	value := t.Value

	env := os.Getenv(value)
	if env == "" {

	}
	return env
}
