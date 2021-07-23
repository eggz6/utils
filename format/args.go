package format

import (
	"regexp"
	"strings"
)

func ReplaceArgs(target string, source map[string]string) string {
	reg := regexp.MustCompile(`(\${[A-Za-z0-9_]+})|(\$[A-Za-z0-9_]+)`)
	args := reg.FindAllString(target, -1)
	for _, arg := range args {
		key := strings.TrimPrefix(arg, "$")
		key = strings.TrimPrefix(key, "{")
		key = strings.TrimSuffix(key, "}")
		val, _ := source[key]
		target = strings.Replace(target, arg, val, 1)
	}

	return target
}
