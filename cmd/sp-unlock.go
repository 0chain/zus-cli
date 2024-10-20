package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

// spUnlock unlocks tokens in stake pool
var spUnlock = &cobra.Command{
	Use:   "sp-unlock",
	Short: "Unlock tokens in stake pool.",
	Long:  `Unlock tokens in stake pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags        = cmd.Flags()
			providerID   string
			providerType sdk.ProviderType
			fee          float64
			err          error
		)

		if flags.Changed("blobber_id") {
			if providerID, err = flags.GetString("blobber_id"); err != nil {
				log.Fatalf("invalid 'blobber_id' flag: %v", err)
			} else {
				providerType = sdk.ProviderBlobber
			}
		} else if flags.Changed("validator_id") {
			if providerID, err = flags.GetString("validator_id"); err != nil {
				log.Fatalf("invalid 'validator_id' flag: %v", err)
			} else {
				providerType = sdk.ProviderValidator
			}
		} else if flags.Changed("authorizer_id") {
			if providerID, err = flags.GetString("authorizer_id"); err != nil {
				log.Fatalf("invalid 'authorizer_id' flag: %v", err)
			} else {
				providerType = sdk.ProviderAuthorizer
			}
		}

		if providerType == 0 || providerID == "" {
			log.Fatal("missing flag: one of 'blobber_id','validator_id' or authorizer_id is required")
		}

		if flags.Changed("fee") {
			if fee, err = flags.GetFloat64("fee"); err != nil {
				log.Fatal("invalid 'fee' flag: ", err)
			}
		}

		unlocked, _, err := sdk.StakePoolUnlock(providerType, providerID, zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatalf("Failed to unlock tokens in stake pool: %v", err)
		}

		// success
		fmt.Printf("tokens unlocked: %d, pool deleted", unlocked)
	},
}

func init() {
	rootCmd.AddCommand(spUnlock)
	spUnlock.PersistentFlags().String("miner_id", "", "for given miner")
	spUnlock.PersistentFlags().String("sharder_id", "", "for given sharder")
	spUnlock.PersistentFlags().String("blobber_id", "", "for given blobber")
	spUnlock.PersistentFlags().String("validator_id", "", "for given validator")
	spUnlock.PersistentFlags().String("authorizer_id", "", "for given authorizer")
	spUnlock.PersistentFlags().Float64("fee", 0.0, "transaction fee, default 0")
}
