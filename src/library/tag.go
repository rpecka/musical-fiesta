package library

import (
	"fiesta/src/csgo/commands"
	"fiesta/src/util"
	"strconv"
	"strings"
)

const unallowedChars = "\n "

func generateTagsFromFilename(trackFilename string) []string {
	trimmed := strings.TrimSpace(trackFilename)
	words := strings.Split(trimmed, " ")
	for index, word := range words {
		word := strings.ToLower(word)
		if !isValidTag(word) {
			continue
		}
		words[index] = word
	}
	uniqueWords := util.Unique(words)
	return uniqueWords
}

func isValidTag(tag string) bool {
	if len(tag) < 1 {
		return false
	}
	if strings.ToLower(tag) != tag {
		return false
	}
	if commands.IsIllegal(tag) {
		return false
	}
	if strings.ContainsAny(tag, unallowedChars) {
		return false
	}
	if _, err := strconv.Atoi(tag); err == nil { // Don't allow integers to avoid collisions with track indices
		return false
	}
	return true
}
