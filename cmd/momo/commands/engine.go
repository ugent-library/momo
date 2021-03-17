package commands

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/storage/es6"
	"github.com/ugent-library/momo/internal/storage/pg"
)

func newEngine() engine.Engine {
	return engine.New(
		engine.WithStore(newStore()),
		engine.WithSearchStore(newSearchStore()),
	)
}

func newStore() engine.Storage {
	store, err := pg.New(viper.GetString("pg-conn"))
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func newSearchStore() engine.SearchStorage {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: strings.Split(viper.GetString("es6-url"), ","),
	})
	if err != nil {
		log.Fatal(err)
	}
	mapping, err := ioutil.ReadFile("etc/es6/rec_mapping.json")
	if err != nil {
		log.Fatal(err)
	}
	store := &es6.Store{
		Client:       client,
		IndexName:    viper.GetString("es6-index"),
		IndexMapping: string(mapping),
	}
	return store
}
