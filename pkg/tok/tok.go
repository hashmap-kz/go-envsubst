package tok

import (
	"fmt"
	"github.com/hashmap.kz/go-envsubst/pkg/cbuf"
	"github.com/hashmap.kz/go-envsubst/pkg/cfg"
	"github.com/hashmap.kz/go-envsubst/pkg/util"
	"log"
	"os"
	"strings"
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

		t, err := nex2(b)
		if err != nil {
			return nil, err
		}

		if t.Type == TokenTypeEof {
			tokens = append(tokens, t)
			break
		}

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
		if err != nil {
			return nil, err
		}

		next, err := b.Nextc()
		if err != nil {
			return nil, err
		}

		// just random text, perhaps commented, anyway: not a var
		if next != '}' {
			unknown := &Token{
				Value: "${" + ident.Value + string(next),
				Type:  TokenTypeTxt,
				Line:  ident.Line,
			}
			return unknown, nil
		}

		ident.Value = fmt.Sprintf("${%s}", ident.Value)
		return ident, err
	}

	// $EDITOR
	if c1 == '$' && util.IsIdentStart(c2) {
		b.Move(1)
		ident, err := getOneIdent(b)
		if err != nil {
			return nil, err
		}

		ident.Value = fmt.Sprintf("$%s", ident.Value)
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
			fmt.Printf("%-7d %-32s %s\n", t.Line, t.Value, os.Getenv(unbraceIdent(t.Value)))
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
		result += fmt.Sprintf("%s", t.Value)
	}
	return result
}

func (tl *Tokenlist) DumpExpanded() string {
	config := cfg.NewConfig()
	result := ""
	for _, t := range tl.Tokens {
		if t == nil {
			break
		}
		if t.Type == TokenTypeEof {
			break
		}
		if t.Type == TokenTypeVar {
			result += expandOneVar(t, config)
		} else {
			result += fmt.Sprintf("%s", t.Value)
		}
	}
	return result
}

func unbraceIdent(ident string) string {
	result := ident
	replace := []string{"$", "{", "}"}
	for _, elem := range replace {
		result = strings.ReplaceAll(result, elem, "")
	}
	return result
}

// TODO: this may be optimized a LOT with the global hashtable

func expandOneVar(t *Token, cfg *cfg.Config) string {
	value := unbraceIdent(t.Value)

	// 1) restricted
	for k := range cfg.Restricted {
		if value == k {
			return t.Value
		}
	}

	// 2) restricted with prefixes
	for k := range cfg.RestrictedWithPrefixes {
		if strings.HasPrefix(value, k) {
			return t.Value
		}
	}

	// get the env-var value itself, or fail
	env := os.Getenv(value)
	if env == "" {
		log.Fatalf("unset variable: %s", value)
	}

	// 3) filters
	if len(cfg.Allowed) > 0 || len(cfg.AllowedWithPrefixes) > 0 {
		for k := range cfg.Allowed {
			if value == k {
				return env
			}
		}
		for k := range cfg.AllowedWithPrefixes {
			if strings.HasPrefix(value, k) {
				return env
			}
		}
		// filters were specified, but:
		// var was not found in specified filters, so:
		// we have to keep it unexpanded
		return t.Value
	}

	// there were no filters, we may expand the var
	return env
}
