package commands

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/ugent-library/momo/records"
)

func init() {
	recIndexCmd.AddCommand(recIndexAllCmd)
	recIndexCmd.AddCommand(recIndexCreateCmd)
	recIndexCmd.AddCommand(recIndexDeleteCmd)

	recCmd.AddCommand(recAddCmd)
	recCmd.AddCommand(recIndexCmd)

	rootCmd.AddCommand(recCmd)
}

var recCmd = &cobra.Command{
	Use:   "rec [command]",
	Short: "Record commands",
}

var recAddCmd = &cobra.Command{
	Use:   "add [file.json ...]",
	Short: "store and index recs",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := records.NewService(newRecordsStore(), newRecordsSearchStore())
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

var recIndexCmd = &cobra.Command{
	Use:   "index [command]",
	Short: "Record index commands",
}

var recIndexAllCmd = &cobra.Command{
	Use:   "all",
	Short: "index all stored records",
	Run: func(cmd *cobra.Command, args []string) {
		service := records.NewService(newRecordsStore(), newRecordsSearchStore())
		err := service.IndexRecs()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var recIndexCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create record search index",
	Run: func(cmd *cobra.Command, args []string) {
		err := newRecordsSearchStore().CreateIndex()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var recIndexDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete record search index",
	Run: func(cmd *cobra.Command, args []string) {
		err := newRecordsSearchStore().DeleteIndex()
		if err != nil {
			log.Fatal(err)
		}
	},
}
