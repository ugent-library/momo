package engine

import (
	"encoding/json"
	"log"
	"os"
)

type LensEngine interface {
	GetLens(name string) *Lens
	Lenses() []*Lens
}

type Lens struct {
	Name  string
	Scope Scope
	Theme string
}

func (e *engine) initLens() {
	lenses := make(map[string]*Lens)
	for _, lens := range loadLenses() {
		lenses[lens.Name] = lens
	}
	e.lenses = lenses
}

func (e *engine) GetLens(name string) *Lens {
	if lens, ok := e.lenses[name]; ok {
		return lens
	}
	return nil
}

func (e *engine) Lenses() []*Lens {
	lenses := make([]*Lens, 0)
	for _, lens := range e.lenses {
		lenses = append(lenses, lens)
	}
	return lenses
}

func loadLenses() []*Lens {
	jsonFile, err := os.Open("etc/lenses.json")
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	lenses := make([]*Lens, 0)
	if err := json.NewDecoder(jsonFile).Decode(&lenses); err != nil {
		log.Fatal(err)
	}
	return lenses
}
