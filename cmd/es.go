package cmd

import (
	"context"

	"github.com/WSBenson/goku/internal/es"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var elasticCmd = &cobra.Command{
	Use:   "es",
	Short: "Serving flexible search capabilities",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		address := viper.GetString("es_address")
		es.ElasticClient(ctx, address)
	},
}

func init() {

	rootCmd.AddCommand(elasticCmd)
	// flag to
	elasticCmd.Flags().StringP("es_address", "a", "localhost:9200", "url for Elasticsearch server")
	elasticCmd.Flags().StringP("es_mapping_file", "m", "./etc/mapping.json", "location for es mapping file")
	// binds the port key to the pflag with viper
	viper.BindPFlag("es_address", elasticCmd.Flags().Lookup("es_address"))
	viper.BindPFlag("es_mapping_file", elasticCmd.Flags().Lookup("es_mapping_file"))
}
