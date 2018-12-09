package chem

import (
	"log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Elements struct {
	Source string `yaml:"source"`
	Elements [] struct {
		Number    string `yaml:"number"`
		Symbol    string `yaml:"symbol"`
		Element   string `yaml:"element"`
		Group     string `yaml:"group"`
		Period    string `yaml:"period"`
		Weight    string `yaml:"weight"`
		Density   string `yaml:"density"`
		Melt      string `yaml:"melt"`
		Boil      string `yaml:"boil"`
		C         string `yaml:"C"`
		X         string `yaml:"X"`
		Abundance string `yaml:"abundance"`
		Property  string `yaml:"property"`
	} `yaml:"elements"`
}

func ReadElements() Elements {
	content, err := ioutil.ReadFile("elements.yaml")
	if err != nil {
		log.Fatal(err)
	}

	e := Elements{}
	err = yaml.Unmarshal([]byte(content), &e)
	if err != nil {
		log.Fatalf("error: %v", err)
	}


	return e
}
