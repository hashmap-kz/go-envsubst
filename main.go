package main

import (
	"flag"
	"fmt"
	"github.com/hashmap.kz/go-envsubst/pkg/tok"
	"github.com/hashmap.kz/go-envsubst/pkg/util"
	"log"
	"strings"
)

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

	b, err := util.ReadFile("data/manifests.yaml")
	if err != nil {
		log.Fatal(err)
	}

	tl, _ := tok.Tokenize(b)
	tl.DumpStat()
	//tl.Dump()
}
