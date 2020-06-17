package cmd

import (
	"github.com/WSBenson/goku/internal/es"
	"github.com/WSBenson/goku/internal/fight"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addfighterCmd = &cobra.Command{
	Use:   "addf",
	Short: "Add your favorite Z fighters with one command",

	Run: func(cmd *cobra.Command, args []string) {
		client := es.NewClient(viper.GetString("addr"), "fighters", viper.GetString("map"))

		f := fight.NewFighter(viper.GetString("name"), viper.GetInt("level"))

		// calls AddF to add a fighter to the fighters index
		//es.CurrClient.AddFighter(f)
		client.AddFighter(f)

		client.QueryFighter(f)
	},
}

func init() {

	rootCmd.AddCommand(addfighterCmd)
	// flag to
	addfighterCmd.Flags().StringP("name", "n", "Goku", "name of fighter to add to elasticsearch")
	addfighterCmd.Flags().IntP("level", "l", 9001, "power level associated with that fighter")
	// binds the port key to the pflag with viper
	viper.BindPFlag("name", addfighterCmd.Flags().Lookup("name"))
	viper.BindPFlag("level", addfighterCmd.Flags().Lookup("level"))
}
