package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

// spLock locks tokens a stake pool lack
var spLock = &cobra.Command{
	Use:   "sp-lock",
	Short: "Lock tokens lacking in stake pool.",
	Long:  `Lock tokens lacking in stake pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags        = cmd.Flags()
			providerID   string
			providerType sdk.ProviderType
			tokens       float64
			fee          float64
			err          error
		)

		if flags.Changed("miner_id") {
			if providerID, err = flags.GetString("miner_id"); err != nil {
				log.Fatalf("invalid 'miner_id' flag: %v", err)
			} else {
				providerType = sdk.ProviderMiner
			}
		} else if flags.Changed("sharder_id") {
			if providerID, err = flags.GetString("sharder_id"); err != nil {
				log.Fatalf("invalid 'sharder_id' flag: %v", err)
			} else {
				providerType = sdk.ProviderSharder
			}
		} else if flags.Changed("blobber_id") {
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
		} else if providerType == 0 || providerID == "" {
			log.Fatal("missing flag: one of 'miner_id', 'sharder_id', 'blobber_id', 'validator_id', 'authorizer_id' is required")
		}

		if !flags.Changed("tokens") {
			log.Fatal("missing required 'tokens' flag")
		}

		if tokens, err = flags.GetFloat64("tokens"); err != nil {
			log.Fatal("invalid 'tokens' flag: ", err)
		}

		if tokens < 0 {
			log.Fatal("invalid token amount: negative")
		}

		if flags.Changed("fee") {
			if fee, err = flags.GetFloat64("fee"); err != nil {
				log.Fatal("invalid 'fee' flag: ", err)
			}
		}

		var hash string
		hash, _, err = sdk.StakePoolLock(providerType, providerID,
			zcncore.ConvertToValue(tokens), zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatalf("Failed to lock tokens in stake pool: %v", err)
		}
		fmt.Println("tokens locked, txn hash:", hash)
	},
}

func init() {
	rootCmd.AddCommand(spLock)
	spLock.PersistentFlags().String("miner_id", "", "for given miner")
	spLock.PersistentFlags().String("sharder_id", "", "for given sharder")
	spLock.PersistentFlags().String("blobber_id", "", "for given blobber")
	spLock.PersistentFlags().String("validator_id", "", "for given validator")
	spLock.PersistentFlags().String("authorizer_id", "", "for given authorizer")
	spLock.PersistentFlags().Float64("tokens", 0.0, "tokens to lock, required")
	spLock.PersistentFlags().Float64("fee", 0.0, "transaction fee, default 0")

	spLock.MarkFlagRequired("tokens")
}
