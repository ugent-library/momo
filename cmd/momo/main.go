package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/cobra"
	"github.com/ugent-library/momo/adding"
	"github.com/ugent-library/momo/http/ui"
	"github.com/ugent-library/momo/storage/es6"
	"github.com/ugent-library/momo/storage/pg"
)

func main() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	mapping, err := ioutil.ReadFile("etc/es6/rec_mapping.json")
	if err != nil {
		log.Fatal(err)
	}
	searchStore := &es6.Store{
		Client:       client,
		IndexName:    "momo_rec",
		IndexMapping: string(mapping),
	}

	rootCmd := &cobra.Command{
		Use:   "momo [command]",
		Short: "The momo CLI",
	}

	var serverPort int
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "The momo webserver",
		Run: func(cmd *cobra.Command, args []string) {
			app := &ui.App{Port: serverPort}
			app.Start()
		},
	}
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 3000, "bind to this TCP port")

	indexCmd := &cobra.Command{
		Use:   "index [command]",
		Short: "Search index operations",
	}
	indexCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create the search index",
		Run: func(cmd *cobra.Command, args []string) {
			err := searchStore.CreateIndex()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	var indexDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the search index",
		Run: func(cmd *cobra.Command, args []string) {
			err := searchStore.DeleteIndex()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	indexCmd.AddCommand(indexCreateCmd)
	indexCmd.AddCommand(indexDeleteCmd)

	recCmd := &cobra.Command{
		Use:   "rec [command]",
		Short: "Rec operations",
	}
	recAddCmd := &cobra.Command{
		Use:   "add [file.json ...]",
		Short: "Add recs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			store, err := pg.New("host=localhost user=nsteenla dbname=momo_dev sslmode=disable")
			if err != nil {
				log.Fatal(err)
			}
			service := adding.NewService(store)
			for _, path := range args {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
				}
				dec := json.NewDecoder(file)
				for {
					var r adding.Rec
					if err := dec.Decode(&r); err == io.EOF {
						break
					} else if err != nil {
						log.Fatal(err)
					}
					service.AddRec(&r)
				}
			}
		},
	}
	recCmd.AddCommand(recAddCmd)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(indexCmd)
	rootCmd.AddCommand(recCmd)
	rootCmd.Execute()
}
