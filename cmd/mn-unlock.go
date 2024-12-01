package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var minerscUnlock = &cobra.Command{
	Use:   "mn-unlock",
	Short: "Unlock miner/sharder stake.",
	Long:  "Unlock miner/sharder stake.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags        = cmd.Flags()
			providerID   string
			providerType zcncore.Provider
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

		_, _, _, _, err = zcncore.MinerSCUnlock(providerID, providerType)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("tokens unlocked")
	},
}

func init() {
	rootCmd.AddCommand(minerscUnlock)
	minerscUnlock.PersistentFlags().String("miner_id", "", "miner ID to lock stake for")
	minerscUnlock.PersistentFlags().String("sharder_id", "", "sharder ID to lock stake for")
}
