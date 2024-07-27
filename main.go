package main

import (
	"fmt"
	cfg2 "github.com/hashmap.kz/go-envsubst/pkg/cfg"
	"github.com/hashmap.kz/go-envsubst/pkg/tok"
	"github.com/hashmap.kz/go-envsubst/pkg/util"
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
	cfg := cfg2.NewConfig()
	fmt.Println(cfg)

	os.Setenv("CI_PROJECT_ROOT_NAMESPACE", "cv")
	os.Setenv("CI_PROJECT_PATH", "cv/system/postgresql")
	os.Setenv("CI_PROJECT_NAME", "postgres")
	os.Setenv("CI_COMMIT_REF_NAME", "dev")
	os.Setenv("APP_IMAGE", "postgres:latest")

	tl := tokenizeFile("data/03-manifests.yaml")
	tl.DumpStat()
	//tl.DumpRawUnexpanded()

	fmt.Println()

	tl = tokenizeFile("data/04-manifests.yaml")
	tl.DumpStat()
	//tl.DumpRawUnexpanded()
}
