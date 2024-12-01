package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var minerscLock = &cobra.Command{
	Use:   "mn-lock",
	Short: "Add miner/sharder stake.",
	Long:  "Add miner/sharder stake.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags        = cmd.Flags()
			providerID   string
			providerType zcncore.Provider
			tokens       float64
			err          error
		)

		if flags.Changed("miner_id") {
			if providerID, err = flags.GetString("miner_id"); err != nil {
				log.Fatalf("invalid 'miner_id' flag: %v", err)
			} else {
				providerType = zcncore.ProviderMiner
			}
		} else if flags.Changed("sharder_id") {
			if providerID, err = flags.GetString("sharder_id"); err != nil {
				log.Fatalf("invalid 'sharder_id' flag: %v", err)
			} else {
				providerType = zcncore.ProviderSharder
			}
		}

		if providerType == 0 || providerID == "" {
			log.Fatal("missing flag: one of 'miner_id' or 'sharder_id' is required")
		}

		if !flags.Changed("tokens") {
			log.Fatal("missing tokens flag")
		}

		if tokens, err = flags.GetFloat64("tokens"); err != nil {
			log.Fatal(err)
		}
		if tokens < 0 {
			log.Fatal("invalid token amount: negative")
		}

		hash, _, _, _, err := zcncore.MinerSCLock(providerID, providerType, zcncore.ConvertToValue(tokens))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("locked with:", hash)
	},
}

func init() {
	rootCmd.AddCommand(minerscLock)
	minerscLock.PersistentFlags().String("miner_id", "", "miner ID to lock stake for")
	minerscLock.PersistentFlags().String("sharder_id", "", "sharder ID to lock stake for")
	minerscLock.PersistentFlags().Float64("tokens", 0, "tokens to lock")
	minerscLock.MarkFlagRequired("tokens")
}
