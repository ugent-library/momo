package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	resetCmd.Flags().Bool("confirm", false, "aborts unless this is set")

	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Destructive reset of the application state",
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("confirm") {
			fmt.Println("Aborting, confirm is not set")
			return
		}
		e := newEngine()
		if err := e.Reset(); err != nil {
			log.Fatal(err)
		}
	},
}
