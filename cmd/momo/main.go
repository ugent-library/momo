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
	defaultPgConn   = "host=localhost dbname=momo_dev sslmode=disable"
	defaultEs6URL   = "http://localhost:9200"
	defaultEs6Index = "momo_rec"
	defaultHost     = "localhost"
	defaultPort     = 3000
)

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

func main() {
	viper.SetEnvPrefix("momo")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd := &cobra.Command{
		Use:   "momo [command]",
		Short: "the momo CLI",
	}
	rootCmd.PersistentFlags().String("pg-conn", defaultPgConn, "PostgreSQL connection string")
	viper.BindPFlag("pg-conn", rootCmd.PersistentFlags().Lookup("pg-conn"))
	viper.SetDefault("pg-conn", defaultPgConn)
	rootCmd.PersistentFlags().String("es6-url", defaultEs6URL, "Elasticsearch 6.x url, separate multiple with comma")
	viper.BindPFlag("es6-url", rootCmd.PersistentFlags().Lookup("es6-url"))
	viper.SetDefault("es6-url", defaultEs6URL)
	rootCmd.PersistentFlags().String("es6-index", defaultEs6Index, "Elasticsearch 6.x index name")
	viper.BindPFlag("es6-index", rootCmd.PersistentFlags().Lookup("es6-index"))
	viper.SetDefault("es6-index", defaultEs6Index)

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
			app.Host = viper.GetString("host")
			app.Port = viper.GetInt("port")
			app.Start()
		},
	}
	serverCmd.Flags().String("host", defaultHost, "server host")
	viper.BindPFlag("host", serverCmd.Flags().Lookup("host"))
	viper.SetDefault("host", defaultHost)
	serverCmd.Flags().Int("port", defaultPort, "server port")
	viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
	viper.SetDefault("port", defaultPort)

	recCmd := &cobra.Command{
		Use:   "rec [command]",
		Short: "rec operations",
	}
	recAddCmd := &cobra.Command{
		Use:   "add [file.json ...]",
		Short: "store and index recs",
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
			time.Sleep(3 * time.Second)
		},
	}
	recIndexCmd := &cobra.Command{
		Use:   "index",
		Short: "index all stored recs",
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
			err = service.IndexRecs()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	recIndexCreateCmd := &cobra.Command{
		Use:   "create-index",
		Short: "create rec search index",
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
		Short: "delete rec search index",
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
	recCmd.AddCommand(recAddCmd)
	recCmd.AddCommand(recIndexCmd)
	recCmd.AddCommand(recIndexCreateCmd)
	recCmd.AddCommand(recIndexDeleteCmd)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(recCmd)
	rootCmd.Execute()
}
