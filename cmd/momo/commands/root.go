package commands

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "momo [command]",
	Short: "The momo CLI",
}

func init() {
	viper.SetEnvPrefix("momo")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.PersistentFlags().String("pg-conn", defaultPgConn, "postgres connection string")
	viper.BindPFlag("pg-conn", rootCmd.PersistentFlags().Lookup("pg-conn"))
	viper.SetDefault("pg-conn", defaultPgConn)
	rootCmd.PersistentFlags().String("es6-url", defaultEs6URL, "elasticsearch 6.x url, separate multiple with comma")
	viper.BindPFlag("es6-url", rootCmd.PersistentFlags().Lookup("es6-url"))
	viper.SetDefault("es6-url", defaultEs6URL)
	rootCmd.PersistentFlags().String("es6-index", defaultEs6Index, "elasticsearch 6.x index name")
	viper.BindPFlag("es6-index", rootCmd.PersistentFlags().Lookup("es6-index"))
	viper.SetDefault("es6-index", defaultEs6Index)
}

// Execute the momo CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
