package cmd

import (
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var lsValidators = &cobra.Command{
	Use:   "ls-validators",
	Short: "Show active Validators.",
	Long:  `Show active Validators in the network.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		doJSON, _ := cmd.Flags().GetBool("json")
		stakable, err := cmd.Flags().GetBool("stakable")
		if err != nil {
			log.Fatalf("error parsing stakable flag: %v", err)
		}

		list, err := sdk.GetValidators(stakable)
		if err != nil {
			log.Fatalf("Failed to get storage SC configurations: %v", err)
		}

		if doJSON {
			util.PrintJSON(list)
		} else {
			util.PrintValidators(list)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsValidators)
	lsValidators.Flags().Bool("json", false, "(default false) pass this flag to get response as json object")
	lsValidators.Flags().Bool("stakable", false, "(default false) Gets only validators that can be staked if set to true")
}
