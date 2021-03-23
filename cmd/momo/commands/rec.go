package commands

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/internal/engine"
)

func init() {
	recGetCmd.Flags().String("format", defaultRecformat, "format")
	viper.BindPFlag("format", recGetCmd.Flags().Lookup("format"))
	viper.SetDefault("format", defaultRecformat)

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
		format := viper.GetString("format")
		encoder := e.NewRecEncoder(os.Stdout, format)
		if encoder == nil {
			log.Fatalf("Unknown format %s", format)
		}

		if len(args) > 0 {
			for _, id := range args {
				rec, err := e.GetRec(id)
				if err != nil {
					log.Fatal(err)
				}
				if err := encoder.Encode(rec); err != nil {
					log.Fatal(err)
				}
			}
		} else {
			c := make(chan *engine.Rec)
			defer close(c)

			go func() {
				for rec := range c {
					if err := encoder.Encode(rec); err != nil {
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
		c := make(chan *engine.Rec)

		var wg sync.WaitGroup

		wg.Add(len(args))

		go func() {
			wg.Wait()
			close(c)
		}()

		addFile := func(f string) {
			defer wg.Done()

			file, err := os.Open(f)
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

				c <- &r
				if verbose {
					p.inc()
				}
			}
		}

		for _, f := range args {
			go addFile(f)
		}

		e.AddRecs(c)

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
