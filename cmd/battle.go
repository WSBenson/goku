package cmd

import (
	"github.com/WSBenson/goku/internal/es"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var battleCmd = &cobra.Command{
	Use:   "battle",
	Short: "Take cover, your Z fighters are about to scrap",

	Run: func(cmd *cobra.Command, args []string) {
		client := es.NewClient(viper.GetString("addr"), "fighters", viper.GetString("map"))

		// calls GetFighters() to query all added fighters
		client.GetFighters()
	},
}

func init() {
	rootCmd.AddCommand(battleCmd)
}
