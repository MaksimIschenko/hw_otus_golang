package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// Regex for proccesing text.
var (
	reSpace = regexp.MustCompile(`\s+`)
	reMarks = regexp.MustCompile(`^[\p{P}\p{S}]*|[\p{P}\p{S}]*$`)
)

// Find 10 most used words per text.
func Top10(text string) []string {
	splitted := PrepareText(text)
	wordMap := make(map[string]int)
	maxCount := 0
	for _, word := range splitted {
		if word == "" || word == "-" {
			continue
		}
		wordMap[word]++
		if wordMap[word] > maxCount {
			maxCount = wordMap[word]
		}
	}
	result := []string{}
	for maxCount > 0 {
		s := []string{}
		for word, count := range wordMap {
			if count == maxCount {
				s = append(s, word)
			}
		}
		sort.Strings(s)
		result = append(result, s...)
		maxCount--
		if len(result) >= 10 {
			break
		}
	}
	if len(result) > 10 {
		return result[:10]
	}
	return result
}

// Prepare text: delete not single spaces and marks before and after word.
func PrepareText(text string) []string {
	replacedText := reSpace.ReplaceAllString(text, " ")
	splttedText := strings.Split(replacedText, " ")
	result := make([]string, 0, len(splttedText))
	for _, word := range splttedText {
		if word == "-" {
			continue
		}
		word = reMarks.ReplaceAllString(word, "")
		if word != "" {
			result = append(result, strings.ToLower(word))
		}
	}
	return result
}
