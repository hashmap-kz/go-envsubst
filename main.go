package main

import (
	"fmt"
	cfg2 "github.com/hashmap.kz/go-envsubst/pkg/cfg"
	"github.com/hashmap.kz/go-envsubst/pkg/tok"
	"github.com/hashmap.kz/go-envsubst/pkg/util"
	"log"
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
	cfg := cfg2.NewConfig()
	fmt.Println(cfg)

	tl := tokenizeFile("data/03-manifests.yaml")
	tl.DumpStat()
	fmt.Println(tl.DumpExpanded())

}
