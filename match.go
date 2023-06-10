package main

import (
	"encoding/json"
	"fmt"
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

type SettingsField struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Refinement string `json:"refinement"`
}

type TitlesToMatch []string

// Fields/Refinement will be unique with possible matches from the title
type FieldToMatches map[SettingsField][]MatchedTitle

type ReturnFieldAndMatches struct {
	Field   SettingsField  `json:"field"`
	Matches []MatchedTitle `json:"matches"`
}

// Custom marshaller to return field and match as JSON
func (fieldMatchings FieldToMatches) MarshalJSON() ([]byte, error) {
	var jsonFieldAndMatches []ReturnFieldAndMatches
	for field, matches := range fieldMatchings {
		jsonFieldAndMatches = append(jsonFieldAndMatches, ReturnFieldAndMatches{
			Field:   field,
			Matches: matches,
		})
	}

	return json.Marshal(jsonFieldAndMatches)
}

// Gets all matching titles for all fields
func FindAllMatches(titles TitlesToMatch) FieldToMatches {
	suggestions := make(FieldToMatches)
	settings := getSettings()
	fields := settings.Category.Fields
	for _, field := range fields {
		for _, tag := range field.Tags {
			for _, title := range titles {

				// Match with label
				score := fuzzy.LevenshteinDistance(
					strings.TrimSpace(strings.ToLower(title)),
					strings.TrimSpace(strings.ToLower(tag.Label)),
				)

				// add if it's close enough
				if score <= MIN_SCORE_TO_MATCH {
					suggestedField := SettingsField{
						Name:       field.Name,
						Type:       field.Type,
						Refinement: tag.Refinement,
					}

					suggestions[suggestedField] = append(suggestions[suggestedField], MatchedTitle{
						OriginalTitle: title,
						Samples:       []string{"test 1", "test 2", "test 3"},
						Score:         score,
					})
				}
			}
		}
	}
	fmt.Printf("%+v\n", suggestions)

	return suggestions
}

// Returns the lowest scored match for each title
func (fieldMatches FieldToMatches) GetBestMatch(ch chan FieldToMatches) {
	fieldBestMatched := make(FieldToMatches)
	var usedTitles []string
	for field, matches := range fieldMatches {
		// Sort suggestions for each field by score
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score < matches[j].Score
		})

		// pick the first (lowest scored) one if it hasn't been used
		if !slices.Contains(usedTitles, matches[0].OriginalTitle) {
			fieldBestMatched[field] = []MatchedTitle{matches[0]}
			usedTitles = append(usedTitles, matches[0].OriginalTitle)
		}
	}
	ch <- fieldBestMatched
}

// Returns the best match for all given field titles
func MatchFields(titles TitlesToMatch, useConcurrency bool) FieldToMatches {
	bestMatches := make(FieldToMatches)
	allMatches := FindAllMatches(titles)
	bestMatchCh := make(chan FieldToMatches)
	numProcesses := 1

	if !useConcurrency {
		go allMatches.GetBestMatch(bestMatchCh)
	}

	if useConcurrency {
		for field, matches := range allMatches {
			newMap := make(FieldToMatches)
			newMap[field] = matches
			go newMap.GetBestMatch(bestMatchCh)
		}
		numProcesses = len(allMatches) * 1
	}

	// Setup concurrency
	for i := 0; i < numProcesses; i++ {
		bestMatch := <-bestMatchCh
		println(bestMatch)
		for k, v := range bestMatch {
			bestMatches[k] = v
		}
	}

	return bestMatches
}
