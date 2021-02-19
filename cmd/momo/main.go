package main

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/spf13/cobra"
	"github.com/ugent-library/momo"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "momo [command]",
		Short: "The momo CLI",
	}

	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "The momo webserver",
		Run: func(cmd *cobra.Command, args []string) {
			app := &momo.App{}
			app.Start()
		},
	}

	var indexCmd = &cobra.Command{
		Use:   "index [command]",
		Short: "Control the search index",
	}
	var indexCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create the search index",
		Run: func(cmd *cobra.Command, args []string) {
			es, err := elasticsearch.NewDefaultClient()
			if err != nil {
				log.Fatalf("Can't create es client: %s", err)
			}
			mappingFile, err := os.Open("config/es/rec_mapping.json")
			defer mappingFile.Close()
			if err != nil {
				log.Fatalf("Can't read mapping file	: %s", err)
			}
			res, err := es.Indices.Create("momo_rec", es.Indices.Create.WithBody(mappingFile))
			if err != nil {
				log.Fatalf("Can't create es index: %s", err)
			}
			if res.IsError() {
				log.Fatalf("Can't create es index: %s", res)
			}

		},
	}
	var indexDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete the search index",
		Run: func(cmd *cobra.Command, args []string) {
			es, err := elasticsearch.NewDefaultClient()
			if err != nil {
				log.Fatalf("Can't create es client: %s", err)
			}
			res, err := es.Indices.Delete([]string{"momo_rec"})
			if err != nil {
				log.Fatalf("Can't delete es index: %s", err)
			}
			if res.IsError() {
				log.Fatalf("Can't delete es index: %s", res)
			}

		},
	}
	indexCmd.AddCommand(indexCreateCmd)
	indexCmd.AddCommand(indexDeleteCmd)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(indexCmd)
	rootCmd.Execute()
}
