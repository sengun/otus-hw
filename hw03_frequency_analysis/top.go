package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var dashExpression = regexp.MustCompile(`(\s-)|(-\s)`)
var wordAppendixExpression = regexp.MustCompile(`[\n\t,.:;"']+`)

type wordCountStruct struct {
	word  string
	count int
}

func Top10(text string) []string {
	counter := make(map[string]int)

	text = PrepareText(text)
	for _, word := range GetTextUnits(text) {
		if len(word) == 0 {
			continue
		}
		lowerCaseWord := strings.ToLower(word)
		counter[lowerCaseWord] = counter[lowerCaseWord] + 1
	}
	wordCounter := make([]wordCountStruct, 0, len(counter))

	for word, wordCount := range counter {
		wordCounter = append(wordCounter, wordCountStruct{word, wordCount})
	}

	return GetFirstTenWords(SortWordCountStruct(wordCounter))
}

func PrepareText(text string) string {
	textWithoutDashes := dashExpression.ReplaceAll([]byte(text), []byte(" "))
	textWithoutAppendix := wordAppendixExpression.ReplaceAll(textWithoutDashes, []byte(" "))

	return string(textWithoutAppendix)
}

func GetTextUnits(text string) []string {
	if text == "" {
		return []string{}
	}

	return strings.Split(text, " ")
}

func GetFirstTenWords(wordCounter []wordCountStruct) []string {
	j := 0
	resultWords := make([]string, 0, len(wordCounter))
	for _, wordCount := range wordCounter {
		resultWords = append(resultWords, wordCount.word)
		j++
		if j >= 10 {
			break
		}
	}

	return resultWords
}

func SortWordCountStruct(wordCounter []wordCountStruct) []wordCountStruct {
	sort.Slice(wordCounter, func(i, j int) bool {
		if wordCounter[i].count == wordCounter[j].count {
			return wordCounter[i].word < wordCounter[j].word
		}
		return wordCounter[i].count > wordCounter[j].count
	})

	return wordCounter
}
