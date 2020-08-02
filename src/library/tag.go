package library

import (
	"fiesta/src/csgo/commands"
	"fiesta/src/csgo/configfile"
	"fiesta/src/util"
	"strconv"
	"strings"
)

const unallowedChars = "\n "

func generateTagsFromFilename(trackFilename string) []string {
	trimmed := strings.TrimSpace(trackFilename)
	words := strings.Split(trimmed, " ")
	validTags := make([]string, 0)
	for _, word := range words {
		word := strings.ToLower(word)
		if !isValidTag(word) {
			continue
		}
		validTags = append(validTags, word)
	}
	uniqueTags := util.Unique(validTags)
	return uniqueTags
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
	if util.Contains(configfile.Commands(), tag) {
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
