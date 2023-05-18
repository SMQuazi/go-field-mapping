package main

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"golang.org/x/exp/slices"
)

var MIN_SCORE_TO_MATCH = 2

type MatchedTitle struct {
	OriginalTitle string   `json:"originalTitle"`
	Samples       []string `json:"samples"`
	Score         int      `json:"score"`
}

type SuggestedField struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Refinement string `json:"refinement"`
}

type TitleForMatching []string

// Fields/Refinement will be unique with possible matches from the title
type FieldsToAllSuggestions map[SuggestedField][]MatchedTitle
type FieldsToOneSuggestion map[SuggestedField]MatchedTitle

type ReturnFieldAndMatch struct {
	Field SuggestedField `json:"field"`
	Match MatchedTitle   `json:"match"`
}

// Custom marshaller to return field and match as JSON
func (fbs FieldsToOneSuggestion) MarshalJSON() ([]byte, error) {
	var fieldsAndMatch []ReturnFieldAndMatch
	for field, match := range fbs {
		fieldsAndMatch = append(fieldsAndMatch, ReturnFieldAndMatch{
			Field: field,
			Match: match,
		})
	}

	return json.Marshal(fieldsAndMatch)
}

// Matches multiple fields for a given title
func SuggestFieldsForOneTitle(header string, ch chan FieldsToAllSuggestions) {
	suggestions := make(FieldsToAllSuggestions)
	settings := getSettings()
	fields := settings.Category.Fields
	for _, field := range fields {
		for _, tag := range field.Tags {
			suggestedField := SuggestedField{
				Name:       field.Name,
				Type:       field.Type,
				Refinement: tag.Refinement,
			}

			// Use refinement if it exists
			score := fuzzy.LevenshteinDistance(
				strings.TrimSpace(strings.ToLower(header)),
				strings.TrimSpace(strings.ToLower(tag.Label)),
			)

			// add if it's close enough
			if score <= MIN_SCORE_TO_MATCH {
				suggestions[suggestedField] = append(suggestions[suggestedField], MatchedTitle{
					OriginalTitle: header,
					Samples:       []string{"test 1", "test 2", "test 3"},
					Score:         score,
				})
			}
		}
	}
	ch <- suggestions
}

// Returns the lowest scored match for each title
func GetBestMatches(allSuggestions FieldsToAllSuggestions, ch chan FieldsToOneSuggestion) {
	fieldBestSuggestion := make(FieldsToOneSuggestion)
	var usedTitles []string
	for field, matches := range allSuggestions {
		// Sort suggestions for each field by score
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score < matches[j].Score
		})

		// pick the first (lowest scored) one if it hasn't been used
		if !slices.Contains(usedTitles, matches[0].OriginalTitle) {
			fieldBestSuggestion[field] = matches[0]
			usedTitles = append(usedTitles, matches[0].OriginalTitle)
		}
	}
	ch <- fieldBestSuggestion
}

// Returns the best match for all given field titles
func MatchFields(titles TitleForMatching) FieldsToOneSuggestion {
	allMatchesChannel := make(chan FieldsToAllSuggestions, len(titles))
	bestMatchChannel := make(chan FieldsToOneSuggestion, len(titles))
	var bestMatches = make(FieldsToOneSuggestion)

	for _, title := range titles {
		go SuggestFieldsForOneTitle(title, allMatchesChannel)
	}

	// Setup concurrency
	numProcesses := len(titles) * 2
	for i := 0; i < numProcesses; i++ {
		select {
		case allMatches := <-allMatchesChannel:

			go GetBestMatches(allMatches, bestMatchChannel)
		case bestMatch := <-bestMatchChannel:
			println(bestMatch)

			// Keep unique list of lowest scores
			for k, v := range bestMatch {
				if _, ok := bestMatches[k]; !ok || bestMatch[k].Score < bestMatches[k].Score {
					bestMatches[k] = v
				}
			}
		}
	}

	return bestMatches
}
