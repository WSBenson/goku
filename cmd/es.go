package cmd

import (
	"github.com/WSBenson/goku/internal/es"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var elasticCmd = &cobra.Command{
	Use:   "es",
	Short: "Serving flexible search capabilities",

	Run: func(cmd *cobra.Command, args []string) {

		address := viper.GetString("addr")
		mapping := viper.GetString("map")
		es.ElasticClient(address, mapping)
	},
}

func init() {

	rootCmd.AddCommand(elasticCmd)
	// flag to
	elasticCmd.Flags().StringP("addr", "a", "http://localhost:9200", "address for Elasticsearch server")
	elasticCmd.Flags().StringP("map", "m", "./mapping/mapping.json", "filepath to mapping.json for fighter index")
	// binds the port key to the pflag with viper
	viper.BindPFlag("addr", elasticCmd.Flags().Lookup("addr"))
	viper.BindPFlag("map", elasticCmd.Flags().Lookup("map"))
}
