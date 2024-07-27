package main

import (
	"fmt"
	"github.com/hashmap.kz/go-envsubst/pkg/tok"
	"github.com/hashmap.kz/go-envsubst/pkg/util"
	"io"
	"log"
	"os"
)

func tokenizeFile(fname string) *tok.Tokenlist {
	b, err := util.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	tl, _ := tok.Tokenize(b)
	return tl
}

func main() {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	content := string(bytes)
	tl, err := tok.Tokenize(content)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(tl.DumpExpanded())
}
