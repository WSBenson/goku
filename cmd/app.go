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
	"github.com/WSBenson/goku/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var port string

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Serving untouchable, instinctive capabilities",

	Run: func(cmd *cobra.Command, args []string) {
		port = viper.GetString("port")
		app.Serve(port)
	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	// flag to
	appCmd.Flags().StringP("port", "p", "3000", "port for goku server")
	// binds the port key to the pflag with viper
	viper.BindPFlag("port", appCmd.Flags().Lookup("port"))
}
