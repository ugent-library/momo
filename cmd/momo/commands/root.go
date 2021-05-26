package commands

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "momo [command]",
	Short: "The momo CLI",
	// flags override env vars
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Changed {
				viper.Set(f.Name, f.Value.String())
			}
		})
		return nil
	},
}

func init() {
	viper.SetEnvPrefix("momo")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("pg-conn", defaultPgConn)
	viper.SetDefault("es6-url", defaultEs6URL)
	viper.SetDefault("es6-index-prefix", defaultEs6IndexPrefix)
	viper.SetDefault("citeproc-url", defaultCiteprocURL)
	viper.SetDefault("base-url", defaultBaseURL)
	viper.SetDefault("port", defaultPort)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().String("pg-conn", defaultPgConn, "postgres connection string")
	rootCmd.PersistentFlags().String("es6-url", defaultEs6URL, "elasticsearch 6.x url, separate multiple with comma")
	rootCmd.PersistentFlags().String("es6-index-prefix", defaultEs6IndexPrefix, "elasticsearch 6.x index prefix")
	rootCmd.PersistentFlags().String("citeproc-url", defaultCiteprocURL, "citeproc url")
}

// Execute the momo CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
