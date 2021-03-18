package engine

import (
	"encoding/json"
	"log"
	"os"
)

type LensEngine interface {
	Lenses() []Lens
}

type Lens struct {
	Name  string
	Scope Scope
	Theme string
}

func (e *engine) initLens() {
	e.lenses = loadLenses()
}

func (e *engine) Lenses() []Lens {
	return e.lenses
}

func loadLenses() []Lens {
	jsonFile, err := os.Open("etc/lenses.json")
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	lenses := make([]Lens, 0)
	if err := json.NewDecoder(jsonFile).Decode(&lenses); err != nil {
		log.Fatal(err)
	}
	return lenses
}
