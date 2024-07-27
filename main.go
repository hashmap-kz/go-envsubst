package main

import (
	"fmt"
	"github.com/hashmap.kz/go-envsubst/pkg/tok"
	"io"
	"log"
	"os"
)

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
