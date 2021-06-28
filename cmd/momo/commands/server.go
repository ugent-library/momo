package commands

import (
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/go-graceful/server"
	"github.com/ugent-library/momo/internal/routes"
)

func init() {
	serverStartCmd.Flags().String("base-url", defaultBaseURL, "base url")
	serverStartCmd.Flags().Int("port", defaultPort, "server port")

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
		r := chi.NewRouter()
		routes.Register(r, newEngine())
		s := server.New(r,
			server.WithPort(viper.GetInt("port")),
		)
		s.Start()
	},
}
