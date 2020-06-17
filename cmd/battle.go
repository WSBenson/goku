package cmd

import (
	"github.com/WSBenson/goku/internal/es"
	"github.com/WSBenson/goku/internal/fight"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var battleCmd = &cobra.Command{
	Use:   "battle",
	Short: "Take cover, your Z fighters are about to scrap",

	Run: func(cmd *cobra.Command, args []string) {
		client := es.NewClient(viper.GetString("addr"), "fighters", viper.GetString("map"))
		f := fight.NewFighter(viper.GetString("name"), viper.GetInt("level"))

		// calls QueryFighter() to query a fighter
		client.QueryFighter(f)
	},
}

func init() {
	rootCmd.AddCommand(battleCmd)

	battleCmd.Flags().StringP("warrior", "f", "Goku", "name of fighter to find with elasticsearch query")
	// binds the port key to the pflag with viper
	viper.BindPFlag("warrior", battleCmd.Flags().Lookup("warrior"))
}
