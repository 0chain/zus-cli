package cmd

import (
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var collectRewards = &cobra.Command{
	Use:   "collect-reward",
	Short: "Collect accrued rewards for a stake pool.",
	Long:  "Collect accrued rewards for a stake pool.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()

		var providerId string
		var err error
		var hash string

		if flags.Changed("provider_id") {
			providerId, err = flags.GetString("provider_id")
			if err != nil {
				log.Fatal(err)
			}
		}

		if !flags.Changed("provider_type") {
			log.Fatal("missing tokens flag")
		}

		providerName, err := flags.GetString("provider_type")
		if err != nil {
			log.Fatal(err)
		}

		switch providerName {
		case "blobber":
			_, _, err = sdk.CollectRewards(providerId, sdk.ProviderBlobber)
		case "validator":
			_, _, err = sdk.CollectRewards(providerId, sdk.ProviderValidator)
		case "miner":
			hash, _, _, _, err = zcncore.MinerSCCollectReward(providerId, zcncore.ProviderMiner)
		case "sharder":
			hash, _, _, _, err = zcncore.MinerSCCollectReward(providerId, zcncore.ProviderSharder)
		case "authorizer":
			hash, _, _, _, err = zcncore.ZCNSCCollectReward(providerId, zcncore.ProviderAuthorizer)
		default:
			log.Fatal("unknown provider type")
		}
		if err != nil {
			log.Fatal("Error paying reward:", err)
		}
		log.Println("transferred reward tokens: ", hash)
	},
}

func init() {
	rootCmd.AddCommand(collectRewards)
	collectRewards.PersistentFlags().String("provider_type", "", "provider type")
	collectRewards.PersistentFlags().String("provider_id", "", "blobber or validator or miner or sharder or authorizer id")
	collectRewards.MarkFlagRequired("provider_id")
	collectRewards.MarkFlagRequired("provider_type")

}
