package main

import (
	"io/ioutil"
	"log"

	"github.com/Universiteitsbibliotheek/momo"
	"github.com/Universiteitsbibliotheek/momo/store"
	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/cobra"
)

func main() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	mapping, err := ioutil.ReadFile("etc/es/rec_mapping.json")
	if err != nil {
		log.Fatal(err)
	}
	esStore := &store.Es{
		Client:       client,
		IndexName:    "momo_rec",
		IndexMapping: string(mapping),
	}

	rootCmd := &cobra.Command{
		Use:   "momo [command]",
		Short: "The momo CLI",
	}

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "The momo webserver",
		Run: func(cmd *cobra.Command, args []string) {
			app := &momo.App{}
			app.Start()
		},
	}

	indexCmd := &cobra.Command{
		Use:   "index [command]",
		Short: "Control the search index",
	}
	indexCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create the search index",
		Run: func(cmd *cobra.Command, args []string) {
			err := esStore.CreateIndex()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	var indexDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the search index",
		Run: func(cmd *cobra.Command, args []string) {
			err := esStore.DeleteIndex()
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
