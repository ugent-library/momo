package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/web"
	"github.com/ugent-library/momo/web/app"
)

func init() {
	serverStartCmd.Flags().String("host", defaultHost, "server host")
	viper.BindPFlag("host", serverStartCmd.Flags().Lookup("host"))
	viper.SetDefault("host", defaultHost)
	serverStartCmd.Flags().Int("port", defaultPort, "server port")
	viper.BindPFlag("port", serverStartCmd.Flags().Lookup("port"))
	viper.SetDefault("port", defaultPort)

	serverCmd.AddCommand(serverStartCmd)

	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server [command]",
	Short: "The momo HTTP server",
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the http server",
	Run: func(cmd *cobra.Command, args []string) {
		a := app.New(newEngine())
		web.RegisterRoutes(a)
		s := web.NewServer(a,
			web.WithHost(viper.GetString("host")),
			web.WithPort(viper.GetInt("port")),
		)
		if err := s.Start(); err != nil {
			log.Fatal(err)
		}
	},
}