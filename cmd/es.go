/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

		// port := viper.GetString("port")
		address := viper.GetString("es_address")
		es.ElasticClient(ctx, address)
	},
}

func init() {

	rootCmd.AddCommand(elasticCmd)
	// flag to
	elasticCmd.Flags().StringP("es_address", "a", "localhost:9200", "url for Elasticsearch server")
	elasticCmd.Flags().StringP("es_mapping_file", "m", "./mapping.json", "location for es mapping file")
	// binds the port key to the pflag with viper
	viper.BindPFlag("es_address", elasticCmd.Flags().Lookup("es_address"))
	viper.BindPFlag("es_mapping_file", elasticCmd.Flags().Lookup("es_mapping_file"))
}
