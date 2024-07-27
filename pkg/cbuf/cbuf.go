package cbuf

import (
	"errors"
	"strings"
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

func (b *CBuf) IsEof() bool {
	return b.Eofs >= EOFS_PADDING_BUFLEN
}

func (b *CBuf) Nextc() (rune, error) {

	// when you build buffer, allocate more space to avoid IOOB check
	// for example: source = { '1', '2', '3', '\0' }, buffer = { '1', '2', '3', '\0', '\0', '\0', '\0', '\0' }

	for {
		if b.IsEof() {
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

func (b *CBuf) Next() string {
	nextc, _ := b.Nextc()
	return string(nextc)
}

func (b *CBuf) Peekc1() (rune, error) {

	// don't be too smart ;)
	saveOffset := b.Offset
	saveLine := b.Line
	saveColumn := b.Column
	savePrevc := b.Prevc
	saveEofs := b.Eofs

	res, _ := b.Nextc()

	b.Offset = saveOffset
	b.Line = saveLine
	b.Column = saveColumn
	b.Prevc = savePrevc
	b.Eofs = saveEofs

	return res, nil
}

func (b *CBuf) Peekc3() ([]rune, error) {

	res := make([]rune, 3)

	// don't be too smart ;)
	saveOffset := b.Offset
	saveLine := b.Line
	saveColumn := b.Column
	savePrevc := b.Prevc
	saveEofs := b.Eofs

	res[0], _ = b.Nextc()
	res[1], _ = b.Nextc()
	res[2], _ = b.Nextc()

	b.Offset = saveOffset
	b.Line = saveLine
	b.Column = saveColumn
	b.Prevc = savePrevc
	b.Eofs = saveEofs

	return res, nil
}

func (b *CBuf) Move(cnt int) {
	for i := 0; i < cnt; i++ {
		_, err := b.Nextc()
		if err != nil {
			return
		}
	}
}
