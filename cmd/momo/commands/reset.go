package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	confirmF bool
)

func init() {
	resetCmd.Flags().BoolVarP(&confirmF, "confirm", "", false, "aborts unless this is set")

	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Destructive reset of the application state",
	Run: func(cmd *cobra.Command, args []string) {
		if !confirmF {
			fmt.Println("Aborting, confirm is not set")
			return
		}
		e := newEngine()
		if err := e.Reset(); err != nil {
			log.Fatal(err)
		}
	},
}
