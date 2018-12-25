package chem

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type element struct {
	Number    int    `yaml:"number"`
	Symbol    string `yaml:"symbol"`
	Element   string `yaml:"element"`
	Group     int    `yaml:"group"`
	Period    int    `yaml:"period"`
	Weight    string `yaml:"weight"`
	Density   string `yaml:"density"`
	Melt      string `yaml:"melt"`
	Boil      string `yaml:"boil"`
	C         string `yaml:"C"`
	X         string `yaml:"X"`
	Abundance string `yaml:"abundance"`
	Property  string `yaml:"property"`
}

type elements struct {
	Source   string     `yaml:"source"`
	Elements [] element `yaml:"elements"`
}

func readElements() elements {
	content, err := ioutil.ReadFile("elements.yaml")
	if err != nil {
		log.Fatal(err)
	}

	e := elements{}
	err = yaml.Unmarshal([]byte(content), &e)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return e
}

type Atom struct {
	Name   string
	Symbol string
	Number int
	Period int
	Group  int
}

type Atoms struct {
	Symbols map[string]Atom
}

func (a *Atoms) read() {
	a.Symbols = make(map[string]Atom)
	elements := readElements()

	for _, e := range elements.Elements {
		atom := Atom{
			Symbol: e.Symbol,
			Name:   e.Element,
			Number: e.Number,
			Period: e.Period,
			Group:  e.Group,
		}
		a.Symbols[e.Symbol] = atom
	}
}

func NewAtoms() *Atoms {
	atom := &Atoms{}
	atom.read()
	return atom
}
