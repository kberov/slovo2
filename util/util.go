// Package util contains utilitiy functions, used across slovo2.
package util

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Slogify converts a sequence of words to lowercased words. Spaces are
// replaced with connector string. Consequitive spaces are treated as one
// space. The default connector is an empty string. If `stripPunct` is true,
// removes any punctuation character.
func Slogify(text string, connector string, stripPunct bool) string {
	words := strings.Split(strings.ToLower(text), " ")
	wordsLen := len(words)

	var slog strings.Builder
	for i, word := range words {
		// Treat consequittive spaces as one by removing empty words, resulted
		// from doubled spaces.
		if len(word) == 0 {
			continue
		}
		var wordSlog strings.Builder
		for _, r := range word {
			if stripPunct && unicode.IsPunct(r) {
				continue
			}
			wordSlog.WriteRune(r)
		}
		slog.WriteString(wordSlog.String())
		if i < wordsLen-1 {
			slog.WriteString(connector)
		}
	}
	return slog.String()
}

// CamelToSnakeCase is used to convert structure fields to
// snake case table columns by sqlx.DB.MapperFunc. See tests for examples.
func CamelToSnakeCase(text string) string {
	if utf8.RuneCountInString(text) == 2 {
		return strings.ToLower(text)
	}
	var snakeCase strings.Builder
	var wordBegins = true
	var prevWasUpper = true
	for _, r := range text {
		wordBegins, prevWasUpper = lowerLetter(&snakeCase, r, wordBegins, prevWasUpper, "_")
	}
	return snakeCase.String()
}

func lowerLetter(snakeCase *strings.Builder, r rune, wordBegins, prevWasUpper bool, connector string) (bool, bool) {
	if unicode.IsUpper(r) && !wordBegins {
		snakeCase.WriteString(connector)
		snakeCase.WriteRune(unicode.ToLower(r))
		wordBegins = true
		prevWasUpper = true
		return wordBegins, prevWasUpper
	}
	// handle case `ID` and beginning of word
	if wordBegins && prevWasUpper {
		snakeCase.WriteRune(unicode.ToLower(r))
		wordBegins = false
		prevWasUpper = false
		return wordBegins, prevWasUpper
	}
	snakeCase.WriteRune(r)
	return wordBegins, prevWasUpper
}
