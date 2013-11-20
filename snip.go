package bluemonday

import (
	"bytes"
	"regexp"
	"strings"
	"unicode/utf8"

	"code.google.com/p/go.net/html"
)

var snipRe = regexp.MustCompile("[\\s]+")

func SnipText(s string, length int) string {
	s = snipRe.ReplaceAllString(strings.TrimSpace(s), " ")
	s = html.UnescapeString(s)
	if len(s) <= length {
		return s
	}
	s = s[:length]
	i := strings.LastIndexAny(s, " .-!?")
	if i != -1 {
		return s[:i]
	}
	return CleanNonUTF8(s)
}

func CleanNonUTF8(s string) string {
	b := &bytes.Buffer{}
	for i := 0; i < len(s); i++ {
		c, size := utf8.DecodeRuneInString(s[i:])
		if c != utf8.RuneError || size != 1 {
			b.WriteRune(c)
		}
	}
	return b.String()
}
