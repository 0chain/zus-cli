package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var updateGlobalConfigCmd = &cobra.Command{
	Use:    "global-update-config",
	Short:  "Update global settings",
	Long:   `Update global settings.`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			hash string
			err  error
		)

		input := new(zcncore.InputMap)
		input.Fields = util.SetupInputMap(cmd.Flags(), "keys", "values")
		if err != nil {
			log.Fatal(err)
		}

		if hash, _, _, _, err = zcncore.MinerScUpdateGlobals(input); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("global settings updated\nHash: %v\n", hash)
	},
}

func init() {
	rootCmd.AddCommand(updateGlobalConfigCmd)
	updateGlobalConfigCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateGlobalConfigCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
