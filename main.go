package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const EOFS_PADDING_BUFLEN = 8
const EOF = -1

type CBuf struct {
	Buffer   string
	Size     int
	OrigSize int
	Offset   int
	Line     int
	Column   int
	Prevc    rune
	Eofs     int
}

func CBufNew(content string) (*CBuf, error) {
	n := content + strings.Repeat(" ", EOFS_PADDING_BUFLEN)

	return &CBuf{
		Buffer:   n,
		Size:     len(n),
		OrigSize: len(content),
		Offset:   0,
		Line:     1,
		Column:   0,
		Prevc:    0,
		Eofs:     -1,
	}, nil
}

func (b *CBuf) isEofIncludePadding() bool {
	return b.Eofs >= EOFS_PADDING_BUFLEN
}

func (b *CBuf) isEof() bool {
	return b.Offset >= b.OrigSize
}

func (b *CBuf) nextc() (rune, error) {

	// when you build buffer, allocate more space to avoid IOOB check
	// for example: source = { '1', '2', '3', '\0' }, buffer = { '1', '2', '3', '\0', '\0', '\0', '\0', '\0' }

	for {
		if b.isEofIncludePadding() {
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

		next, _ := utf8.DecodeRuneInString(b.Buffer[b.Offset:])
		b.Offset += 1
		b.Column += 1
		b.Prevc = next

		if next == 0 {
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

func main() {
	b, _ := CBufNew("some\\\nstring")
	for {
		if b.isEof() {
			break
		}
		peek3, _ := b.peekc3()
		fmt.Println(string(peek3))
		b.nextc()
	}
}
