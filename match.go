package main

import (
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type MatchSuggestions struct {
	suggestedField Field
	originalTitle  string
	score          int
}

func scoreMatch(titles []string) (suggestions []MatchSuggestions) {
	settings := getSettings()
	fields := settings.Category.Fields
	fieldMap := make(map[string][]Field)
	for _, field := range fields {
		for _, title := range titles {
			score := fuzzy.LevenshteinDistance(strings.ToLower(title), strings.ToLower(field.Name))
			suggestions = append(suggestions, MatchSuggestions{
				originalTitle:  title,
				suggestedField: field,
				score:          score,
			})
		}
	}
	return
}
