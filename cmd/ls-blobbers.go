package cmd

import (
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

// lsBlobers shows active blobbers
var lsBlobers = &cobra.Command{
	Use:   "ls-blobbers",
	Short: "Show active blobbers in storage SC.",
	Long:  `Show active blobbers in storage SC.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		doJSON, _ := cmd.Flags().GetBool("json")
		doAll, _ := cmd.Flags().GetBool("all")
		isStakable, err := cmd.Flags().GetBool("stakable")
		if err != nil {
			log.Fatalf("err parsing in stakable flag: %v", err)
		}
		// set is_active=true to get only active blobbers
		isActive := true
		if doAll {
			isActive = false
		}
		list, err := sdk.GetBlobbers(isActive, isStakable)
		if err != nil {
			log.Fatalf("Failed to get storage SC configurations: %v", err)
		}

		if doJSON {
			util.PrintJSON(list)
		} else {
			util.PrintBlobbers(list, isActive)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsBlobers)
	lsBlobers.Flags().Bool("json", false, "(default false) pass this option to print response as json data")
	lsBlobers.Flags().Bool("all", false, "(default false) shows active and non active list of blobbers on ls-blobbers")
	lsBlobers.Flags().Bool("stakable", false, "(default false) gets only stakable list of blobbers if set to true")
}
