package cmd

import (
	"github.com/WSBenson/goku/internal/kafku"
	"github.com/spf13/cobra"
)

var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Test kafka",

	Run: func(cmd *cobra.Command, args []string) {
		kafku.ProduceZ("goku")
	},
}

func init() {

	rootCmd.AddCommand(kafkaCmd)

}
