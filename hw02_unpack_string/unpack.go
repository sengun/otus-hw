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
		case i+2 < len(runes) && isSymbolEscapedWithCount(runes[i], runes[i+1], runes[i+2]):
			count, _ := strconv.Atoi(string(runes[i+2]))
			output.WriteString(strings.Repeat(string(runes[i+1]), count))

			i += 3

		case i+1 < len(runes) && isEscapeRune(runes[i]) && isEscapeRuneOrDigit(runes[i+1]):
			output.WriteString(strings.Repeat(string(runes[i+1]), 1))

			i += 2

		case i+1 < len(runes) && isSymbolNotEscapedWithCount(runes[i], runes[i+1]):
			count, _ := strconv.Atoi(string(runes[i+1]))
			output.WriteString(strings.Repeat(string(runes[i]), count))

			i += 2

		case !isEscapeRuneOrDigit(runes[i]):
			output.WriteString(string(runes[i]))

			i++

		default:
			return "", ErrInvalidString
		}
	}

	return output.String(), nil
}

func isSymbolEscapedWithCount(rune1 rune, rune2 rune, rune3 rune) bool {
	if !isEscapeRune(rune1) {
		return false
	}

	if !isEscapeRuneOrDigit(rune2) {
		return false
	}

	if !unicode.IsDigit(rune3) {
		return false
	}

	return true
}

func isEscapeRune(symbol rune) bool {
	return symbol == escapeRune
}

func isSymbolNotEscapedWithCount(rune1 rune, rune2 rune) bool {
	if isEscapeRuneOrDigit(rune1) {
		return false
	}

	if !unicode.IsDigit(rune2) {
		return false
	}

	return true
}

func isEscapeRuneOrDigit(symbol rune) bool {
	return isEscapeRune(symbol) || unicode.IsDigit(symbol)
}
