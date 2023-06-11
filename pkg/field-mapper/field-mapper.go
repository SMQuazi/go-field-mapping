package fieldmapper

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

type MappedTitle struct {
	OriginalTitle string   `json:"originalTitle"`
	Samples       []string `json:"samples"`
	Score         int      `json:"score"`
}

type MappedField struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Refinement string `json:"refinement"`
}

type Titles []string

// Fields/Refinement will be unique with possible matches from the title
type MappedFieldsAndTitles map[MappedField][]MappedTitle

type JsonMappedFieldsAndTitles struct {
	Field   MappedField   `json:"field"`
	Matches []MappedTitle `json:"matches"`
}

// Custom marshaller to return field and match as JSON
func (fieldMatchings MappedFieldsAndTitles) MarshalJSON() ([]byte, error) {
	var jsonFieldAndMatches []JsonMappedFieldsAndTitles
	for field, matches := range fieldMatchings {
		jsonFieldAndMatches = append(jsonFieldAndMatches, JsonMappedFieldsAndTitles{
			Field:   field,
			Matches: matches,
		})
	}

	return json.Marshal(jsonFieldAndMatches)
}

// Gets all matching titles for all fields
func FindAllMatches(titles Titles, pathToSettings string) MappedFieldsAndTitles {
	suggestions := make(MappedFieldsAndTitles)
	settings := getSettings(pathToSettings)
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
					suggestedField := MappedField{
						Name:       field.Name,
						Type:       field.Type,
						Refinement: tag.Refinement,
					}

					suggestions[suggestedField] = append(suggestions[suggestedField], MappedTitle{
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
func (fieldMatches MappedFieldsAndTitles) GetBestMatch(ch chan MappedFieldsAndTitles) {
	fieldBestMatched := make(MappedFieldsAndTitles)
	var usedTitles []string
	for field, matches := range fieldMatches {
		// Sort suggestions for each field by score
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score < matches[j].Score
		})

		// pick the first (lowest scored) one if it hasn't been used
		if !slices.Contains(usedTitles, matches[0].OriginalTitle) {
			fieldBestMatched[field] = []MappedTitle{matches[0]}
			usedTitles = append(usedTitles, matches[0].OriginalTitle)
		}
	}
	ch <- fieldBestMatched
}

// Returns the best match for all given field titles
func MatchFields(titles Titles, useConcurrency bool, pathToSettings string) MappedFieldsAndTitles {
	bestMatches := make(MappedFieldsAndTitles)
	allMatches := FindAllMatches(titles, pathToSettings)
	bestMatchCh := make(chan MappedFieldsAndTitles)
	numProcesses := 1

	if !useConcurrency {
		go allMatches.GetBestMatch(bestMatchCh)
	}

	if useConcurrency {
		for field, matches := range allMatches {
			newMap := make(MappedFieldsAndTitles)
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

	// return non matches as custom
	for _, title := range titles {
		titleFound := false
		for _, matches := range bestMatches {
			if matches[0].OriginalTitle == title {
				titleFound = true
				break
			}
		}

		if !titleFound {
			customField := MappedField{
				Name:       "custom_field",
				Type:       "str",
				Refinement: "",
			}
			bestMatches[customField] = append(bestMatches[customField], MappedTitle{
				OriginalTitle: title,
				Samples:       []string{"test 1", "test 2", "test 3"},
				Score:         UNMATCHED_SCORE,
			})
		}
	}

	return bestMatches
}
