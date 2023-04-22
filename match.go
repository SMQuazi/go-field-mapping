package main

import "github.com/lithammer/fuzzysearch/fuzzy"

type MatchSuggestions struct {
	originalTitle  string
	suggestedField Field
	score          int
}

func findMatches(titles []string) (suggestions []MatchSuggestions) {
	settings := getSettings()
	fields := settings.Category.Fields[0]
	for _, title := range titles {
		score := fuzzy.RankMatch(title, fields.Name)
		currentField := fields
		suggestions = append(suggestions, MatchSuggestions{
			originalTitle:  title,
			suggestedField: currentField,
			score:          score,
		})
	}
	return
}
