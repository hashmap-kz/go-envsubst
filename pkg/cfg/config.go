package cfg

import (
	"log"
	"os"
	"strings"
)

const (
	GenvsubstAllowed                = "GENVSUBST_ALLOWED"
	GenvsubstAllowedWithPrefixes    = "GENVSUBST_ALLOWED_WITH_PREFIXES"
	GenvsubstRestricted             = "GENVSUBST_RESTRICTED"
	GenvsubstRestrictedWithPrefixes = "GENVSUBST_RESTRICTED_WITH_PREFIXES"
)

type Config struct {
	Allowed                map[string]string
	AllowedWithPrefixes    map[string]string
	Restricted             map[string]string
	RestrictedWithPrefixes map[string]string
}

func NewConfig() *Config {

	config := &Config{
		Allowed:                parseList(GenvsubstAllowed),
		AllowedWithPrefixes:    parseList(GenvsubstAllowedWithPrefixes),
		Restricted:             parseList(GenvsubstRestricted),
		RestrictedWithPrefixes: parseList(GenvsubstRestrictedWithPrefixes),
	}

	return config
}

func parseList(env string) map[string]string {
	result := map[string]string{}

	elem := strings.TrimSpace(os.Getenv(env))
	if elem == "" {
		return result
	}

	if strings.Contains(elem, " ") && strings.Contains(elem, ",") {
		log.Fatalf("cannot use both spaces and commas in setting: %s\n", env)
	}

	if strings.Contains(elem, " ") {
		for _, s := range strings.Split(elem, " ") {
			s = strings.TrimSpace(s)
			result[s] = s
		}
	} else if strings.Contains(elem, ",") {
		for _, s := range strings.Split(elem, ",") {
			s = strings.TrimSpace(s)
			result[s] = s
		}
	} else {
		// just single value
		result[elem] = elem
	}

	return result
}
