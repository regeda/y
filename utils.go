package y

import (
	"regexp"
	"strings"
)

// Values is the map of something
type Values map[string]interface{}

var (
	snake = regexp.MustCompile("([A-Z]*)([A-Z][^A-Z]+|$)")
)

func underscore(s string) string {
	var a []string
	for _, sub := range snake.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, "_"))
}
