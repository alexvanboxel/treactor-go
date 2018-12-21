package chem

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type element struct {
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
	Symbols string
}

type Atoms struct {
	Symbols map[string]Atom
}

func (a *Atoms) read() {
	a.Symbols = make(map[string]Atom)
	elements := readElements()
	fmt.Println(elements)

	for _, e := range elements.Elements {
		a.Symbols[e.Symbol] = Atom{
			Symbols: e.Symbol,
		}
	}
}

func NewAtoms() *Atoms {
	atom := &Atoms{}
	atom.read()
	return atom
}
