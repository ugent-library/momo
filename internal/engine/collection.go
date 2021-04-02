package engine

import (
	"encoding/json"
	"log"
	"os"
)

type CollectionEngine interface {
	Collections() []Collection
}

type Collection struct {
	Name  string
	Theme string
}

func (e *engine) initCollections() {
	e.collections = loadCollections()
}

func (e *engine) Collections() []Collection {
	return e.collections
}

func loadCollections() []Collection {
	jsonFile, err := os.Open("etc/collections.json")
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	lenses := make([]Collection, 0)
	if err := json.NewDecoder(jsonFile).Decode(&lenses); err != nil {
		log.Fatal(err)
	}
	return lenses
}
