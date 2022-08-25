package redact

import (
	"fmt"
	"strings"
)

// Word returns the first character of the word, followed by "*"s for each addtional character
func Word(word string) string {
	if word == "" {
		return ""
	}

	runes := []rune(word)
	l := len(runes) - 1
	firstChar := string(runes[0:1])
	return fmt.Sprintf("%s%s", firstChar, strings.Repeat("*", l))
}

// Words redacts each word in a string.
func Words(words string) string {
	if words == "" {
		return ""
	}

	ws := strings.Split(words, " ")
	for i := range ws {
		ws[i] = Word(ws[i])
	}
	return strings.Join(ws, " ")
}

// Email redacts the part of an email before the @.
func Email(email string) string {
	if email == "" {
		return ""
	}

	parts := strings.Split(email, "@")
	parts[0] = Words(parts[0])
	return strings.Join(parts, "@")
}

// Phone redacts the last digits of a phone number.
func Phone(phone string) string {
	if phone == "" {
		return ""
	}

	parts := strings.Split(phone, "-")
	lastRunes := []rune(parts[len(parts)-1])
	lastPartlen := len(lastRunes)
	repeatChars := 4
	if repeatChars > lastPartlen {
		repeatChars = lastPartlen
	}

	parts[len(parts)-1] = string(append(lastRunes[0:lastPartlen-repeatChars], []rune(strings.Repeat("*", repeatChars))...))
	return strings.Join(parts, "-")
}
