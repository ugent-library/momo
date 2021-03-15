package commands

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/ugent-library/momo/engine"
)

func init() {
	recIndexCmd.AddCommand(recIndexAllCmd)
	recIndexCmd.AddCommand(recIndexCreateCmd)
	recIndexCmd.AddCommand(recIndexDeleteCmd)

	recCmd.AddCommand(recGetCmd)
	recCmd.AddCommand(recAddCmd)
	recCmd.AddCommand(recIndexCmd)

	rootCmd.AddCommand(recCmd)
}

var recCmd = &cobra.Command{
	Use:   "rec [command]",
	Short: "Record commands",
}

var recGetCmd = &cobra.Command{
	Use:   "get [id ...]",
	Short: "get stored records",
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		enc := json.NewEncoder(os.Stdout)

		if len(args) > 0 {
			for _, id := range args {
				rec, err := e.GetRec(id)
				if err != nil {
					log.Fatal(err)
				}
				if err := enc.Encode(rec); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			c := make(chan *engine.Rec)
			defer close(c)

			go func() {
				for rec := range c {
					if err := enc.Encode(rec); err != nil {
						log.Fatal(err)
					}
				}
			}()

			if err := e.GetAllRecs(c); err != nil {
				log.Fatal(err)
			}
		}
	},
}

var recAddCmd = &cobra.Command{
	Use:   "add [file.json ...]",
	Short: "store and index recs",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		p := newProgress(100)
		out := make(chan *engine.Rec)
		e.AddRecs(out)

		// TODO read json files concurrently?
		for _, path := range args {
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			dec := json.NewDecoder(file)
			for {
				var r engine.Rec
				if err := dec.Decode(&r); err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}
				out <- &r
				if verbose {
					p.inc()
				}
			}
		}

		close(out)

		if verbose {
			p.done()
		}
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
		e := newEngine()
		err := e.IndexRecs()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var recIndexCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create record search index",
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		err := e.CreateRecIndex()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var recIndexDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete record search index",
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		err := e.DeleteRecIndex()
		if err != nil {
			log.Fatal(err)
		}
	},
}
