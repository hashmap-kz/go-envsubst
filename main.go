package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const EOFS_PADDING_BUFLEN = 8
const EOF = -1

type CBuf struct {
	Buffer string
	Size   int
	Offset int
	Line   int
	Column int
	Prevc  rune
	Eofs   int
}

func CBufNew(content string) (*CBuf, error) {
	n := content + strings.Repeat(" ", EOFS_PADDING_BUFLEN)

	return &CBuf{
		Buffer: n,
		Size:   len(content),
		Offset: 0,
		Line:   1,
		Column: 0,
		Prevc:  0,
		Eofs:   -1,
	}, nil
}

func (b *CBuf) isEof() bool {
	return b.Eofs >= EOFS_PADDING_BUFLEN
}

func (b *CBuf) nextc() (rune, error) {

	// when you build buffer, allocate more space to avoid IOOB check
	// for example: source = { '1', '2', '3', '\0' }, buffer = { '1', '2', '3', '\0', '\0', '\0', '\0', '\0' }

	for {
		if b.isEof() {
			break
		}

		if b.Eofs > EOFS_PADDING_BUFLEN {
			return EOF, errors.New("infinite loop handling")
		}

		if b.Prevc == '\n' {
			b.Line += 1
			b.Column = 0
		}

		if b.Buffer[b.Offset] == '\\' {
			if b.Buffer[b.Offset+1] == '\r' {
				if b.Buffer[b.Offset+2] == '\n' {
					// DOS: [\][\r][\n]
					b.Offset += 3
				} else {
					// OSX: [\][\r]
					b.Offset += 2
				}

				b.Prevc = '\n'
				continue
			}

			// UNX: [\][\n]
			if b.Buffer[b.Offset+1] == '\n' {
				b.Offset += 2
				b.Prevc = '\n'
				continue
			}
		}

		if b.Buffer[b.Offset] == '\r' {
			if b.Buffer[b.Offset+1] == '\n' {
				// DOS: [\r][\n]
				b.Offset += 2
			} else {
				// OSX: [\r]
				b.Offset += 1
			}
			b.Prevc = '\n'
			return '\n', nil
		}

		if b.Offset == b.Size {
			b.Eofs += 1
			return EOF, nil
		}

		next, cnt := utf8.DecodeRuneInString(b.Buffer[b.Offset:])
		b.Offset += 1
		b.Column += 1
		b.Prevc = next

		if cnt == 0 || next == utf8.RuneError {
			b.Eofs += 1
			return EOF, nil
		}

		return next, nil
	}

	return EOF, nil
}

func (b *CBuf) next() string {
	nextc, _ := b.nextc()
	return string(nextc)
}

func (b *CBuf) peekc1() (rune, error) {

	// don't be too smart ;)
	saveOffset := b.Offset
	saveLine := b.Line
	saveColumn := b.Column
	savePrevc := b.Prevc
	saveEofs := b.Eofs

	res, _ := b.nextc()

	b.Offset = saveOffset
	b.Line = saveLine
	b.Column = saveColumn
	b.Prevc = savePrevc
	b.Eofs = saveEofs

	return res, nil
}

func (b *CBuf) peekc3() ([]rune, error) {

	res := make([]rune, 3)

	// don't be too smart ;)
	saveOffset := b.Offset
	saveLine := b.Line
	saveColumn := b.Column
	savePrevc := b.Prevc
	saveEofs := b.Eofs

	res[0], _ = b.nextc()
	res[1], _ = b.nextc()
	res[2], _ = b.nextc()

	b.Offset = saveOffset
	b.Line = saveLine
	b.Column = saveColumn
	b.Prevc = savePrevc
	b.Eofs = saveEofs

	return res, nil
}

func (b *CBuf) move(cnt int) {
	for i := 0; i < cnt; i++ {
		_, err := b.nextc()
		if err != nil {
			return
		}
	}
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentTail(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

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
	b, err := CBufNew(content)
	if err != nil {
		return nil, err
	}

	tokens := []*Token{}
	for {
		if b.isEof() {
			break
		}
		t, _ := nex2(b)
		tokens = append(tokens, t)
	}

	return &Tokenlist{
		Tokens: tokens,
	}, nil
}

func nex2(b *CBuf) (*Token, error) {
	peek3, _ := b.peekc3()
	c1 := peek3[0]
	c2 := peek3[1]
	c3 := peek3[2]

	if c1 == '$' && c2 == '{' && isIdentStart(c3) {
		b.move(2)
		value := ""
		for {
			peek, err := b.peekc1()
			if err != nil {
				return nil, err
			}
			if !isIdentTail(peek) {
				break
			}

			next := b.next()
			value += next
		}
		next, _ := b.nextc()
		if next != '}' {
			return nil, errors.New("unclosed_brace")
		}
		return &Token{
			Value: value,
			Type:  TokenTypeVar,
			Line:  b.Line,
		}, nil
	}

	if b.isEof() {
		return &Token{
			Value: "",
			Type:  TokenTypeEof,
			Line:  b.Line,
		}, nil
	}

	n, err := b.nextc()
	if err != nil {
		return nil, err
	}
	if n == EOF {
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

func readFile(name string) (string, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (tl *Tokenlist) dumpStat() {
	for _, t := range tl.Tokens {
		if t == nil {
			break
		}
		if t.Type == TokenTypeEof {
			break
		}
		if t.Type == TokenTypeVar {
			fmt.Printf("var at line: %d, named: %s\n", t.Line, t.Value)
		}
	}
}

func (tl *Tokenlist) dump() {
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

// TODO: clean

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var allowedVars arrayFlags
var forbiddenVars arrayFlags
var allowedWithPrefixVars arrayFlags
var forbiddenWithPrefixVars arrayFlags

func parseListFlags(where arrayFlags) map[string]string {
	tmp := map[string]string{}
	for _, elem := range where {
		elem = strings.TrimSpace(elem)
		if strings.Contains(elem, " ") && strings.Contains(elem, ",") {
			log.Fatal("cannot use both spaces and commas in flags: " + elem)
		}

		if strings.Contains(elem, " ") {
			for _, s := range strings.Split(elem, " ") {
				tmp[s] = s
			}
		}
		if strings.Contains(elem, ",") {
			for _, s := range strings.Split(elem, ",") {
				tmp[s] = s
			}
		}
	}
	return tmp
}

func main() {
	// "${var} no var"

	flag.Var(&allowedVars, "allowed", "Expand only allowed, ignore others")
	flag.Var(&forbiddenVars, "forbidden", "Never expand these vars, this flag has the highest priority")

	flag.Var(&allowedWithPrefixVars, "allowedWithPrefix", "Expand only allowed, ignore others")
	flag.Var(&forbiddenWithPrefixVars, "forbiddenWithPrefix", "Never expand these vars, this flag has the highest priority")
	flag.Parse()

	fmt.Println(parseListFlags(allowedVars))
	fmt.Println(parseListFlags(forbiddenVars))

	b, err := readFile("data/manifests.yaml")
	if err != nil {
		log.Fatal(err)
	}

	tl, _ := Tokenize(b)
	tl.dumpStat()
	//tl.dump()
}
