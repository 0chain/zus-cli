package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

// rpUnlock unlocks tokens in a read pool
var rpUnlock = &cobra.Command{
	Use:   "rp-unlock",
	Short: "Unlock some expired tokens in a read pool.",
	Long:  `Unlock some expired tokens in a read pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags = cmd.Flags()
			fee   float64
			err   error
		)

		if flags.Changed("fee") {
			if fee, err = flags.GetFloat64("fee"); err != nil {
				log.Fatal("invalid 'fee' flag: ", err)
			}
		}

		_, _, err = sdk.ReadPoolUnlock(zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatalf("Failed to unlock tokens in read pool: %v", err)
		}
		fmt.Println("unlocked")
	},
}

func init() {
	rootCmd.AddCommand(rpUnlock)
	rpUnlock.PersistentFlags().Float64("fee", 0.0,
		"transaction fee, default 0")
}
