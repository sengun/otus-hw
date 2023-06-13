package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const escapeRune = '\\'

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	runes := []rune(input)
	output := strings.Builder{}

	for i := 0; i < len(runes); {
		switch {
		case i+2 < len(runes) && isSymbolEscapedWithCount(runes[i], runes[i+2]):
			count, _ := strconv.Atoi(string(runes[i+2]))
			output.WriteString(strings.Repeat(string(runes[i+1]), count))

			i += 3

		case i+1 < len(runes) && isSymbolEscapedWithoutCount(runes[i]):
			output.WriteString(strings.Repeat(string(runes[i+1]), 1))

			i += 2

		case i+1 < len(runes) && isSymbolNotEscapedWithCount(runes[i], runes[i+1]):
			count, _ := strconv.Atoi(string(runes[i+1]))
			output.WriteString(strings.Repeat(string(runes[i]), count))

			i += 2

		case !unicode.IsDigit(runes[i]) && runes[i] != escapeRune:
			output.WriteString(string(runes[i]))

			i++

		default:
			return "", ErrInvalidString
		}
	}

	return output.String(), nil
}

func isSymbolEscapedWithCount(rune1 rune, rune3 rune) bool {
	if rune1 != escapeRune {
		return false
	}

	if !unicode.IsDigit(rune3) {
		return false
	}

	return true
}

func isSymbolEscapedWithoutCount(rune1 rune) bool {
	return rune1 == escapeRune
}

func isSymbolNotEscapedWithCount(rune1 rune, rune2 rune) bool {
	if rune1 == escapeRune || unicode.IsDigit(rune1) {
		return false
	}

	if !unicode.IsDigit(rune2) {
		return false
	}

	return true
}
