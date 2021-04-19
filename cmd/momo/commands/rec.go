package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/formats/csljson"
)

var (
	formatF string
	queryF  string
)

func init() {
	recGetCmd.Flags().StringVarP(&formatF, "format", "f", defaultRecformat, "export format")

	recSearchCmd.Flags().StringVarP(&formatF, "format", "f", defaultRecformat, "export format")
	recSearchCmd.Flags().StringVarP(&queryF, "query", "q", "", "search query")

	recIndexCmd.AddCommand(recIndexAllCmd)
	recIndexCmd.AddCommand(recIndexCreateCmd)
	recIndexCmd.AddCommand(recIndexDeleteCmd)

	recCmd.AddCommand(recGetCmd)
	recCmd.AddCommand(recSearchCmd)
	recCmd.AddCommand(recAddCmd)
	recCmd.AddCommand(recIndexCmd)
	recCmd.AddCommand(recAddCitationsCmd)

	rootCmd.AddCommand(recCmd)
}

var recCmd = &cobra.Command{
	Use:   "rec [command]",
	Short: "Record commands",
}

var recGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get stored records",
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		encoder := e.NewRecEncoder(os.Stdout, formatF)
		if encoder == nil {
			log.Fatalf("Unknown format %s", formatF)
		}

		e.EachRec(func(rec *engine.Rec) bool {
			if err := encoder.Encode(rec); err != nil {
				log.Fatal(err)
			}
			return true
		})
	},
}

var recSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "search records",
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		encoder := e.NewRecEncoder(os.Stdout, formatF)
		if encoder == nil {
			log.Fatalf("Unknown format %s", formatF)
		}

		e.SearchEachRec(engine.SearchArgs{Query: queryF}, func(rec *engine.Rec) bool {
			if err := encoder.Encode(rec); err != nil {
				log.Fatal(err)
			}
			return true
		})
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

var recAddCitationsCmd = &cobra.Command{
	Use:   "add-citations",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		p := newProgress(100)

		client := http.Client{
			Timeout: 10 * time.Second,
		}

		e.EachRec(func(rec *engine.Rec) bool {
			body := struct {
				Items []json.RawMessage `json:"items"`
			}{}
			var buf bytes.Buffer
			encoder := csljson.NewEncoder(&buf)
			if err := encoder.Encode(rec); err != nil {
				log.Fatal(err)
			}
			body.Items = append(body.Items, buf.Bytes())
			jsonBody, _ := json.Marshal(body)
			req, err := http.NewRequest("POST", "http://127.0.0.1:8085?style=mla", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			cites, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			cite := gjson.GetBytes(cites, "bibliography.1.0").String()

			if len(cite) == 0 {
				return true
			}

			rep := &engine.Representation{
				RecID:  rec.ID,
				Format: "mla",
				Data:   []byte(cite),
			}

			if err = e.AddRepresentation(rep); err != nil {
				log.Print(err)
			}

			if verbose {
				p.inc()
			}

			return true
		})

		if verbose {
			p.done()
		}
	},
}
