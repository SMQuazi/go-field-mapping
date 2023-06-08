package main

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"golang.org/x/exp/slices"
)

var MIN_SCORE_TO_MATCH = 3
var UNMATCHED_SCORE = 100

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

type TitlesForMatching []string

// Fields/Refinement will be unique with possible matches from the title
type FieldToAllSuggestions map[SuggestedField][]MatchedTitle
type FieldToSuggestion map[SuggestedField]MatchedTitle

type ReturnFieldAndMatch struct {
	Field SuggestedField `json:"field"`
	Match MatchedTitle   `json:"match"`
}

// Custom marshaller to return field and match as JSON
func (fbs FieldToSuggestion) MarshalJSON() ([]byte, error) {
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
func SuggestFieldsForOneTitle(header string, ch chan FieldToAllSuggestions) {
	suggestions := make(FieldToAllSuggestions)
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
func (allSuggestions FieldToAllSuggestions) GetBestMatches(ch chan FieldToSuggestion) {
	fieldBestSuggestion := make(FieldToSuggestion)
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
func MatchFields(titles TitlesForMatching) FieldToSuggestion {
	allMatchesChannel := make(chan FieldToAllSuggestions, len(titles))
	bestMatchChannel := make(chan FieldToSuggestion, len(titles))
	var bestMatches = make(FieldToSuggestion)

	for _, title := range titles {
		go SuggestFieldsForOneTitle(title, allMatchesChannel)
	}

	// Setup concurrency
	numProcesses := len(titles) * 2
	for i := 0; i < numProcesses; i++ {
		select {
		case allMatches := <-allMatchesChannel:
			go allMatches.GetBestMatches(bestMatchChannel)
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

	// map unmapped titles to custom_types
	customField := SuggestedField{
		Name:       "custom_field",
		Type:       "str",
		Refinement: "",
	}
	for _, title := range titles {
		titleFound := false
		for _, matchedTitle := range bestMatches {
			if matchedTitle.OriginalTitle == title {
				titleFound = true
				break
			}
		}
		if !titleFound {
			bestMatches[customField] = MatchedTitle{
				OriginalTitle: title,
				Samples:       []string{""},
				Score:         UNMATCHED_SCORE,
			}
		}

	}

	return bestMatches
}
