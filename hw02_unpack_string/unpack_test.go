package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `a3\a3`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestIsSymbolEscapedWithCount(t *testing.T) {
	dataSet := []struct {
		rune1          rune
		rune2          rune
		rune3          rune
		expectedResult bool
	}{
		{rune1: '\\', rune2: '\\', rune3: '3', expectedResult: true},
		{rune1: '\\', rune2: '3', rune3: '4', expectedResult: true},
		{rune1: '\\', rune2: 'a', rune3: '1', expectedResult: false},
		{rune1: '\\', rune2: '3', rune3: '\\', expectedResult: false},
		{rune1: '3', rune2: 'a', rune3: '1', expectedResult: false},
	}

	for _, ds := range dataSet {
		ds := ds
		t.Run(string(ds.rune1)+string(ds.rune2)+string(ds.rune3), func(t *testing.T) {
			result := isSymbolEscapedWithCount(ds.rune1, ds.rune2, ds.rune3)
			require.Equal(t, ds.expectedResult, result)
		})
	}
}

func TestIsSymbolNotEscapedWithCount(t *testing.T) {
	dataSet := []struct {
		rune1          rune
		rune2          rune
		expectedResult bool
	}{
		{rune1: 'a', rune2: '3', expectedResult: true},
		{rune1: 'a', rune2: '\\', expectedResult: false},
		{rune1: '\\', rune2: 'a', expectedResult: false},
		{rune1: '\\', rune2: '3', expectedResult: false},
		{rune1: '3', rune2: 'a', expectedResult: false},
	}

	for _, ds := range dataSet {
		ds := ds
		t.Run(string(ds.rune1)+string(ds.rune2), func(t *testing.T) {
			result := isSymbolNotEscapedWithCount(ds.rune1, ds.rune2)
			require.Equal(t, ds.expectedResult, result)
		})
	}
}

func TestIsEscapeRune(t *testing.T) {
	require.Equal(t, true, isEscapeRune('\\'))
	require.Equal(t, false, isEscapeRune('3'))
	require.Equal(t, false, isEscapeRune('a'))
}

func TestIsEscapeRuneOrDigit(t *testing.T) {
	dataSet := []struct {
		symbol         rune
		expectedResult bool
	}{
		{symbol: 'a', expectedResult: false},
		{symbol: 'z', expectedResult: false},
		{symbol: '\\', expectedResult: true},
		{symbol: '3', expectedResult: true},
	}

	for _, ds := range dataSet {
		ds := ds
		t.Run(string(ds.symbol), func(t *testing.T) {
			result := isEscapeRuneOrDigit(ds.symbol)
			require.Equal(t, ds.expectedResult, result)
		})
	}
}
