package cmd

import (
	"log"

	"github.com/0chain/gosdk/core/transaction"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

// scConfig shows SC configurations
var scConfig = &cobra.Command{
	Use:   "sc-config",
	Short: "Show storage SC configuration.",
	Long:  `Show storage SC configuration.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		doJSON, _ := cmd.Flags().GetBool("json")

		var conf, err = transaction.GetConfig("storage_sc_config")
		if err != nil {
			log.Fatalf("Failed to get storage SC configurations: %v", err)
		}
		if doJSON {
			util.PrintJSON(conf)
			return
		}
		util.PrintMap(conf.Fields)
	},
}

func init() {
	rootCmd.AddCommand(scConfig)
	scConfig.Flags().Bool("json", false, "(default false) pass this option to print response as json data")
}
