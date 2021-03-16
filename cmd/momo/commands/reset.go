package commands

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Destructive reset of the application state",
	Run: func(cmd *cobra.Command, args []string) {
		if !confirm() {
			return
		}
		e := newEngine()
		if err := e.Reset(); err != nil {
			log.Fatal(err)
		}
	},
}
