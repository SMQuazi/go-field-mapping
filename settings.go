package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Alternatives struct {
	Label      string `yaml:"label"`
	Refinement string `yaml:"refinement,omitempty"`
}

type Fields struct {
	Name         string         `yaml:"name"`
	Type         string         `yaml:"type"`
	Alternatives []Alternatives `yaml:"alternatives"`
}

type Category struct {
	Fields []Fields `yaml:"fields"`
}

type Settings struct {
	Category Category `yaml:"category"`
}

func getSettings() (Settings, error) {
	data, err := os.ReadFile("settings.yml")
	if err != nil {
		var set Settings
		return set, err
	}
	var settings Settings
	ymlErr := yaml.Unmarshal(data, &settings)
	if ymlErr != nil {
		return settings, nil
	}
	return settings, nil
}
