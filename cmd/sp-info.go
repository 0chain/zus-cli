package cmd

import (
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

// spInfo information
var spInfo = &cobra.Command{
	Use:   "sp-info",
	Short: "Stake pool information.",
	Long:  `Stake pool information.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags        = cmd.Flags()
			err          error
			providerID   string
			providerType sdk.ProviderType
		)

		doJSON, _ := cmd.Flags().GetBool("json")

		if flags.Changed("blobber_id") {
			if providerID, err = flags.GetString("blobber_id"); err != nil {
				log.Fatalf("Error: cannot get the value of blobber_id")
			} else {
				providerType = sdk.ProviderBlobber
			}
		} else if flags.Changed("validator_id") {
			if providerID, err = flags.GetString("validator_id"); err != nil {
				log.Fatalf("Error: cannot get the value of validator_id")
			} else {
				providerType = sdk.ProviderValidator
			}
		} else if flags.Changed("authorizer_id") {
			if providerID, err = flags.GetString("authorizer_id"); err != nil {
				log.Fatalf("Error: cannot get the value of authorizer_id")
			} else {
				providerType = sdk.ProviderAuthorizer
			}
		}

		if providerType == 0 || providerID == "" {
			log.Fatal("Error: missing flag: one of 'blobber_id','validator_id' or authorizer_id is required")
		}

		var info *sdk.StakePoolInfo
		if info, err = sdk.GetStakePoolInfo(providerType, providerID); err != nil {
			log.Fatalf("Failed to get stake pool info: %v", err)
		}
		if doJSON {
			util.PrintJSON(info)
		} else {
			util.PrintStakePoolInfo(info)
		}
	},
}

func init() {
	rootCmd.AddCommand(spInfo)
	spInfo.PersistentFlags().String("miner_id", "", "for given miner")
	spInfo.PersistentFlags().String("sharder_id", "", "for given sharder")
	spInfo.PersistentFlags().String("blobber_id", "", "for given blobber")
	spInfo.PersistentFlags().String("validator_id", "", "for given validator")
	spInfo.PersistentFlags().String("authorizer_id", "", "for given authorizer")
	spInfo.PersistentFlags().Bool("json", false, "(default false) pass this option to print response as json data")
}
