package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/http/ui"
	"github.com/ugent-library/momo/records"
	"github.com/ugent-library/momo/storage/es6"
	"github.com/ugent-library/momo/storage/pg"
)

const (
	pgConnDefault = "host=localhost dbname=momo_dev sslmode=disable"
	portDefault   = 3000
)

func newRecordsStore() (records.Storage, error) {
	store, err := pg.New(viper.GetString("pg-conn"))
	if err != nil {
		return nil, err
	}
	return store, nil
}

func newRecordsSearchStore() (records.SearchStorage, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	mapping, err := ioutil.ReadFile("etc/es6/rec_mapping.json")
	if err != nil {
		return nil, err
	}
	store := &es6.Store{
		Client:       client,
		IndexName:    "momo_rec",
		IndexMapping: string(mapping),
	}
	return store, nil
}

func main() {
	viper.SetEnvPrefix("momo")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd := &cobra.Command{
		Use:   "momo [command]",
		Short: "the momo CLI",
	}
	rootCmd.PersistentFlags().String("pg-conn", pgConnDefault, "PostgreSQL connection string")
	viper.BindPFlag("pg-conn", rootCmd.PersistentFlags().Lookup("pg-conn"))
	viper.SetDefault("pg-conn", pgConnDefault)

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "the momo HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			store, err := newRecordsStore()
			if err != nil {
				log.Fatal(err)
			}
			searchStore, err := newRecordsSearchStore()
			if err != nil {
				log.Fatal(err)
			}
			app := ui.New(store, searchStore)
			app.Port = viper.GetInt("port")
			app.Start()
		},
	}
	serverCmd.Flags().Int("port", portDefault, "server port")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
	viper.SetDefault("port", portDefault)

	recCmd := &cobra.Command{
		Use:   "rec [command]",
		Short: "Rec operations",
	}
	recAddCmd := &cobra.Command{
		Use:   "add [file.json ...]",
		Short: "Add recs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			store, err := newRecordsStore()
			if err != nil {
				log.Fatal(err)
			}
			searchStore, err := newRecordsSearchStore()
			if err != nil {
				log.Fatal(err)
			}
			service := records.NewService(store, searchStore)
			out := make(chan *records.Rec)
			service.AddRecs(out)

			// parse json files
			// TODO read json files concurrently?
			for _, path := range args {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
				}
				dec := json.NewDecoder(file)
				for {
					var r records.Rec
					if err := dec.Decode(&r); err == io.EOF {
						break
					} else if err != nil {
						log.Fatal(err)
					}
					out <- &r
				}
			}

			close(out)

			// TODO flush stdio or send output back over channel?
			time.Sleep(5 * time.Second)
		},
	}

	recIndexCreateCmd := &cobra.Command{
		Use:   "create-index",
		Short: "Create the search index",
		Run: func(cmd *cobra.Command, args []string) {
			store, err := newRecordsSearchStore()
			if err != nil {
				log.Fatal(err)
			}
			err = store.CreateIndex()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	recIndexDeleteCmd := &cobra.Command{
		Use:   "delete-index",
		Short: "Delete the search index",
		Run: func(cmd *cobra.Command, args []string) {
			store, err := newRecordsSearchStore()
			if err != nil {
				log.Fatal(err)
			}
			err = store.DeleteIndex()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	recCmd.AddCommand(recIndexCreateCmd)
	recCmd.AddCommand(recIndexDeleteCmd)
	recCmd.AddCommand(recAddCmd)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(recCmd)
	rootCmd.Execute()
}
