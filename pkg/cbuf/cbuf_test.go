package cbuf

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestBufferEmpty(t *testing.T) {

	b, err := CBufNew("")
	assert.NoError(t, err)
	assert.Equal(t, 0, b.Size)

	nextc, err := b.Nextc()
	assert.NoError(t, err)
	assert.Equal(t, EOF, nextc)

}

func mustGetc(b *CBuf) string {
	nextc, err := b.Nextc()
	if err != nil {
		log.Fatal(err)
	}
	return string(nextc)
}

func TestBufferFull(t *testing.T) {
	content := `ca\
fe`

	assert.Equal(t, "ca\\\nfe", content)

	b, err := CBufNew(content)
	assert.NoError(t, err)

	assert.Equal(t, 6, b.Size)

	assert.Equal(t, "c", mustGetc(b))
	assert.Equal(t, "a", mustGetc(b))
	assert.Equal(t, "f", mustGetc(b))
	assert.Equal(t, "e", mustGetc(b))
}

func TestBufferPeek(t *testing.T) {
	b, err := CBufNew("abcde")
	assert.NoError(t, err)

	peekc3, err := b.Peekc3()
	assert.NoError(t, err)
	assert.Equal(t, []rune{'a', 'b', 'c'}, peekc3)

	assert.Equal(t, 0, b.Offset)
	assert.Equal(t, "a", mustGetc(b))
	assert.Equal(t, 1, b.Offset)
}
