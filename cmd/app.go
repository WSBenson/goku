package cmd

import (
	"github.com/WSBenson/goku/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Serving untouchable, instinctive capabilities",

	Run: func(cmd *cobra.Command, args []string) {
		port := viper.GetString("port")
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
