package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/momo/internal/sitemap"
)

func init() {
	sitemapCmd.Flags().String("base-url", defaultBaseURL, "base url")

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
		if err := sitemap.Generate(e, viper.GetString("base-url")); err != nil {
			log.Fatal(err)
		}
	},
}
