package chem

import (
	"github.com/alexvanboxel/reactor/pkg/config"
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
	Symbol  string
	Service string
	Port    string
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
		}
		if config.Mode == "local" {
			atom.Service = "localhost"
			atom.Port = config.Port
		} else {
			log.Fatal("Only local is supported")
		}

		a.Symbols[e.Symbol] = atom
	}
}

func NewAtoms() *Atoms {
	atom := &Atoms{}
	atom.read()
	return atom
}
