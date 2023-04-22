package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Alternatives struct {
	Label      string `yaml:"label"`
	Refinement string `yaml:"refinement,omitempty"`
}

type Field struct {
	Name         string         `yaml:"name"`
	Type         string         `yaml:"type"`
	Alternatives []Alternatives `yaml:"alternatives"`
}

type Category struct {
	Fields []Field `yaml:"fields"`
}

type Settings struct {
	Category Category `yaml:"category"`
}

func getSettings() Settings {
	data, err := os.ReadFile("settings.yml")
	if err != nil {
		log.Fatal(err)
	}
	var settings Settings
	ymlErr := yaml.Unmarshal(data, &settings)
	if ymlErr != nil {
		log.Fatal(ymlErr)
	}
	return settings
}
