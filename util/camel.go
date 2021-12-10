package util

import (
	"bytes"
	"strings"
)

// Just suit for camel, snake, kebab
func ConvertString2UpperCamel(s string) string {
	// split words
	words := make([]string, 0)
	buf := bytes.Buffer{}
	for i := range s {
		// for snake and kebab
		if s[i] == '_' || s[i] == '-' {
			if buf.Len() > 0 {
				words = append(words, word2UpperCamel(buf.String()))
				buf.Reset()
			}
			continue
		}

		// for camel
		if 'A' <= s[i] && s[i] <= 'Z' {
			if buf.Len() > 0 {
				words = append(words, word2UpperCamel(buf.String()))
				buf.Reset()
			}
			buf.WriteByte(s[i])
			continue
		}
		buf.WriteByte(s[i])
	}

	if buf.Len() > 0 {
		words = append(words, word2UpperCamel(buf.String()))
	}

	return strings.Join(words, "")
}

func word2UpperCamel(word string) string {
	if len(word) == 0 {
		return ""
	}
	s := []byte(strings.ToLower(word))
	// A: 0100 0001
	// a: 0110 0001
	s[0] &= 0b11011111
	return string(s)
}
