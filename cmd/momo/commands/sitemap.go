package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/internal/sitemap"
)

func init() {
	sitemapGenerateCmd.Flags().String("host", defaultHost, "server host")
	sitemapGenerateCmd.Flags().Int("port", defaultPort, "server port")

	sitemapCmd.AddCommand(sitemapGenerateCmd)
	rootCmd.AddCommand(sitemapCmd)
}

var sitemapCmd = &cobra.Command{
	Use:   "sitemap [command]",
	Short: "Sitemap commands",
}

var sitemapGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate a new sitemap",
	Run: func(cmd *cobra.Command, args []string) {
		e := newEngine()
		url := fmt.Sprintf("http://%s:%d", viper.GetString("host"), viper.GetInt("port"))
		if err := sitemap.Generate(e, url); err != nil {
			log.Fatal(err)
		}
	},
}
