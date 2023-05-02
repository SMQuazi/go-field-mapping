package main

import (
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var MIN_SCORE_TO_MATCH = 2

type SuggestedMatch struct {
	OriginalTitle string   `json:"originalTitle"`
	TargetField   Field    `json:"targetField"`
	Refinement    string   `json:"refinement"`
	Samples       []string `json:"samples"`
	Score         int      `json:"score"`
}

type SuggestedField struct {
	Name       string `json:"Name"`
	Type       string `json:"Type"`
	Label      string `json:"Label"`
	Refinement string `json:"Refinement"`
}

type FieldsAllSuggestions map[SuggestedField][]SuggestedMatch
type FieldsBestSuggestion map[SuggestedField]SuggestedMatch

func (fbs FieldsBestSuggestion) Unmarshal(bytes []byte) error {
	//TODO create custom marshaller to return JSON
}

func suggestFieldsForTitles(headers []string) FieldsBestSuggestion {
	suggestions := make(FieldsAllSuggestions)
	settings := getSettings()
	fields := settings.Category.Fields
	for _, field := range fields {
		for _, tag := range field.Tags {
			for _, header := range headers {
				suggestedField := SuggestedField{
					Name:       field.Name,
					Type:       field.Type,
					Label:      tag.Label,
					Refinement: tag.Refinement,
				}

				// Use refinement if it exists
				var labelOrRefinement string
				if len(tag.Refinement) > 0 {
					labelOrRefinement = tag.Refinement
				} else {
					labelOrRefinement = tag.Label
				}
				score := fuzzy.LevenshteinDistance(
					strings.TrimSpace(strings.ToLower(header)),
					strings.TrimSpace(strings.ToLower(labelOrRefinement)),
				)

				// add if it's close enough
				if score <= MIN_SCORE_TO_MATCH {
					suggestions[suggestedField] = append(suggestions[suggestedField], SuggestedMatch{
						OriginalTitle: header,
						TargetField:   field,
						Refinement:    tag.Refinement,
						Samples:       []string{"test1", "test2", "test 3"},
						Score:         score,
					})
				}
			}
		}
	}
	bestSuggestions := suggestions.pickBestMatch()
	return bestSuggestions
}

func (allSuggestions FieldsAllSuggestions) pickBestMatch() FieldsBestSuggestion {
	fieldBestSuggestion := make(FieldsBestSuggestion)
	for allSuggestionsField, allSuggestionsMatches := range allSuggestions {
		sort.Slice(allSuggestionsMatches, func(i, j int) bool {
			return allSuggestionsMatches[i].Score < allSuggestionsMatches[j].Score
		})
		fieldBestSuggestion[allSuggestionsField] = allSuggestions[allSuggestionsField][0]
	}
	return fieldBestSuggestion
}
