package cmd

import (
	"log"

	"github.com/0chain/gosdk/core/transaction"
	"github.com/0chain/zus-cli/util"

	"github.com/spf13/cobra"
)

var mnConfigCmd = &cobra.Command{
	Use:   "mn-config",
	Short: "Show miner SC configuration.",
	Long:  `Show miner SC configuration.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			fields = new(transaction.InputMap)
			err    error
		)
		if fields, err = transaction.GetConfig("miner_sc_configs"); err != nil {
			log.Fatal(err)
		}

		util.PrintMap(fields.Fields)
	},
}

func init() {
	rootCmd.AddCommand(mnConfigCmd)
}
