package fieldmapper

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Tag struct {
	Label      string `yaml:"label"`
	Refinement string `yaml:"refinement"`
}

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Tags []Tag  `yaml:"tags"`
}

type Category struct {
	Fields []Field `yaml:"fields"`
}

type Settings struct {
	Category Category `yaml:"category"`
}

func getSettings(path string) Settings {
	data, err := os.ReadFile(path)
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
