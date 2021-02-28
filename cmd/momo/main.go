package main

import (
	"io/ioutil"
	"log"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/cobra"
	"github.com/ugent-library/momo/storage/es6"
	"github.com/ugent-library/momo/ui"
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
	store := &es6.Store{
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
		Short: "Control the search index",
	}
	indexCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create the search index",
		Run: func(cmd *cobra.Command, args []string) {
			err := store.CreateIndex()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	var indexDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the search index",
		Run: func(cmd *cobra.Command, args []string) {
			err := store.DeleteIndex()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	indexCmd.AddCommand(indexCreateCmd)
	indexCmd.AddCommand(indexDeleteCmd)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(indexCmd)
	rootCmd.Execute()
}
