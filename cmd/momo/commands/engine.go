package commands

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/formats/jsonl"
	"github.com/ugent-library/momo/internal/formats/oaidc"
	"github.com/ugent-library/momo/internal/formats/ris"
	"github.com/ugent-library/momo/internal/i18n/gettext"
	"github.com/ugent-library/momo/internal/storage/es6"
	"github.com/ugent-library/momo/internal/storage/sql"
)

func newEngine() engine.Engine {
	return engine.New(
		engine.WithStore(newStore()),
		engine.WithSearchStore(newSearchStore()),
		engine.WithRecEncoder("json", jsonl.NewEncoder),
		engine.WithRecEncoder("oai_dc", oaidc.NewEncoder),
		engine.WithRecEncoder("ris", ris.NewEncoder),
		engine.WithI18n(gettext.New()),
	)
}

func newStore() engine.Storage {
	store, err := sql.New("postgres", viper.GetString("pg-conn"))
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func newSearchStore() engine.SearchStorage {
	mapping, err := ioutil.ReadFile("etc/es6/rec_mapping.json")
	if err != nil {
		log.Fatal(err)
	}
	store, err := es6.New(es6.Config{
		ClientConfig: elasticsearch.Config{
			Addresses: strings.Split(viper.GetString("es6-url"), ","),
		},
		IndexPrefix:  viper.GetString("es6-index-prefix"),
		IndexMapping: string(mapping),
	})
	if err != nil {
		log.Fatal(err)
	}
	return store
}
