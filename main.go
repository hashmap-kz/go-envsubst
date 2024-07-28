package main

import (
	"flag"
	"fmt"
	"github.com/hashmap.kz/go-envsubst/pkg/tok"
	"io"
	"log"
	"os"
)

var (
	dryRun = flag.Bool("dry-run", false, "")
)

func main() {
	flag.Parse()

	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	content := string(bytes)
	tl, err := tok.Tokenize(content)
	if err != nil {
		log.Fatal(err)
	}

	if *dryRun {
		tl.DumpStat()
	} else {
		fmt.Print(tl.DumpExpanded())
	}
}
