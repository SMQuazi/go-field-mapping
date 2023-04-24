package main

import (
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type ScoredTitle struct {
	originalTitle string
	score         int
}

func scoreTitles(titles []string) (suggestions map[*Field][]ScoredTitle) {
	settings := getSettings()
	fields := settings.Category.Fields
	for _, field := range fields {
		for _, title := range titles {
			score := fuzzy.LevenshteinDistance(strings.ToLower(title), strings.ToLower(field.Name))

			// if no records exist for the field in the suggestions, add it
			_, ok := suggestions[&field]
			if !ok {
				suggestions[&field][0] = ScoredTitle{
					originalTitle: title,
					score:         score,
				}
				continue
			}

			// if there is, loop through them, add
			for index, scoredTitle := range suggestions[&field] {
				if scoredTitle.score < score {
					suggestions[&field][index] = ScoredTitle{
						originalTitle: title,
						score:         score,
					}
				}
			}

		}
	}
	return
}
