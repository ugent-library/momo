package commands

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/records"
	"github.com/ugent-library/momo/storage/es6"
	"github.com/ugent-library/momo/storage/pg"
)

var rootCmd = &cobra.Command{
	Use:   "momo [command]",
	Short: "The momo CLI",
}

func init() {
	viper.SetEnvPrefix("momo")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().String("pg-conn", defaultPgConn, "postgres connection string")
	viper.BindPFlag("pg-conn", rootCmd.PersistentFlags().Lookup("pg-conn"))
	viper.SetDefault("pg-conn", defaultPgConn)
	rootCmd.PersistentFlags().String("es6-url", defaultEs6URL, "elasticsearch 6.x url, separate multiple with comma")
	viper.BindPFlag("es6-url", rootCmd.PersistentFlags().Lookup("es6-url"))
	viper.SetDefault("es6-url", defaultEs6URL)
	rootCmd.PersistentFlags().String("es6-index", defaultEs6Index, "elasticsearch 6.x index name")
	viper.BindPFlag("es6-index", rootCmd.PersistentFlags().Lookup("es6-index"))
	viper.SetDefault("es6-index", defaultEs6Index)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func newRecordsStore() (records.Storage, error) {
	store, err := pg.New(viper.GetString("pg-conn"))
	if err != nil {
		return nil, err
	}
	return store, nil
}

func newRecordsSearchStore() (records.SearchStorage, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: strings.Split(viper.GetString("es6-url"), ","),
	})
	if err != nil {
		return nil, err
	}
	mapping, err := ioutil.ReadFile("etc/es6/rec_mapping.json")
	if err != nil {
		return nil, err
	}
	store := &es6.Store{
		Client:       client,
		IndexName:    viper.GetString("es6-index"),
		IndexMapping: string(mapping),
	}
	return store, nil
}
