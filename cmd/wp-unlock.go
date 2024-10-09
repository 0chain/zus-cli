package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

// wpUnlock unlocks tokens in a write pool
var wpUnlock = &cobra.Command{
	Use:   "wp-unlock",
	Short: "Unlock some expired tokens in a write pool.",
	Long:  `Unlock some expired tokens in a write pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags   = cmd.Flags()
			allocID string
			fee     float64
			err     error
		)

		if !flags.Changed("allocation") {
			log.Fatal("missing required 'allocation' flag")
		}
		if allocID, err = flags.GetString("allocation"); err != nil {
			log.Fatal("invalid 'allocation' flag: ", err)
		}

		if flags.Changed("fee") {
			if fee, err = flags.GetFloat64("fee"); err != nil {
				log.Fatal("invalid 'fee' flag: ", err)
			}
		}

		_, _, err = sdk.WritePoolUnlock(allocID, zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatalf("Failed to unlock tokens in write pool: %v", err)
		}
		fmt.Println("unlocked")
	},
}

func init() {
	rootCmd.AddCommand(wpUnlock)
	wpUnlock.PersistentFlags().String("allocation", "",
		"allocation id from which to unlock tokens")
	wpUnlock.PersistentFlags().Float64("fee", 0.0,
		"transaction fee, default 0")

	wpUnlock.MarkFlagRequired("allocation")
}
