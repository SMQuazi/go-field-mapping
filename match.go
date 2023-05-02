package main

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var MIN_SCORE_TO_MATCH = 2

type SuggestedMatch struct {
	OriginalTitle string   `json:"originalTitle"`
	Refinement    string   `json:"refinement"`
	Samples       []string `json:"samples"`
	Score         int      `json:"score"`
}

type SuggestedField struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Label      string `json:"label"`
	Refinement string `json:"refinement"`
}

type TitlesToMatchFields []string
type FieldsAllSuggestions map[SuggestedField][]SuggestedMatch
type FieldsBestSuggestion map[SuggestedField]SuggestedMatch

type ReturnFieldAndMatch struct {
	Field SuggestedField `json:"field"`
	Match SuggestedMatch `json:"match"`
}

func (fbs FieldsBestSuggestion) MarshalJSON() ([]byte, error) {
	var fieldsAndMatch []ReturnFieldAndMatch
	for suggestedField, suggestedMatch := range fbs {
		fieldsAndMatch = append(fieldsAndMatch, ReturnFieldAndMatch{
			Field: suggestedField,
			Match: suggestedMatch,
		})
	}

	return json.Marshal(fieldsAndMatch)
}

func (headers TitlesToMatchFields) SuggestFieldsForTitles() FieldsAllSuggestions {
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
						Refinement:    tag.Refinement,
						Samples:       []string{"test1", "test2", "test 3"},
						Score:         score,
					})
				}
			}
		}
	}
	return suggestions
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
