// Package util contains utilitiy functions, used across slovo2.
package util

import (
	"strings"
	"unicode"
)

// Slogify converts a title or upper cased sequence of words to lowercased
// words. Spaces are replaced with connector string. The default
// connector is an empty string. If `stripPunct` is true, removes any
// punctuation character.
func Slogify(text string, connector string, stripPunct bool) string {
	var slog strings.Builder
	words := strings.Split(text, " ")
	wordsLen := len(words)

	for i, word := range words {
		word = strings.ToLower(word)
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

// ToSnakeCase is used to convert structure fields to
// snake case table columns by sqlx.DB.MapperFunc. See tests for examples.
func ToSnakeCase(text string) string {
	var snakeCase strings.Builder
	var wordBoundary = true
	var prevWasUpper = true
	for _, r := range text {
		wordBoundary, prevWasUpper = lowerLetter(&snakeCase, r, wordBoundary, prevWasUpper, "_")
	}
	return snakeCase.String()
}

func lowerLetter(snakeCase *strings.Builder, r rune, wordBoundary, prevWasUpper bool, connector string) (bool, bool) {
	if unicode.IsUpper(r) && !wordBoundary {
		snakeCase.WriteString(connector)
		snakeCase.WriteRune(unicode.ToLower(r))
		wordBoundary = true
		prevWasUpper = true
		return wordBoundary, prevWasUpper
	}
	// handle case `ID` and beginning of word
	if wordBoundary && prevWasUpper {
		snakeCase.WriteRune(unicode.ToLower(r))
		wordBoundary = false
		prevWasUpper = false
		return wordBoundary, prevWasUpper
	}
	snakeCase.WriteRune(r)
	return wordBoundary, prevWasUpper
}
