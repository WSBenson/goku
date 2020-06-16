package cmd

import (
	"context"

	"github.com/WSBenson/goku/internal"
	"github.com/WSBenson/goku/internal/es"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addfighterCmd = &cobra.Command{
	Use:   "addf",
	Short: "Add your favorite Z fighters with one command",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client, err := elastic.NewSimpleClient(elastic.SetURL(viper.GetString("addr")))
		if err != nil {
			// Handle error
			internal.Logger.Fatal().Err(err).Msg("failed to make new elastic search client")
		}

		// calls AddF to add a fighter to the fighters index
		es.AddF(ctx, client)
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
