package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

// rpLock locks tokens in read pool
var rpLock = &cobra.Command{
	Use:   "rp-lock",
	Short: "Lock some tokens in read pool.",
	Long:  `Lock some tokens in read pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags  = cmd.Flags()
			tokens float64
			fee    float64
			err    error
		)

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

		_, _, err = sdk.ReadPoolLock(zcncore.ConvertToValue(tokens), zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatalf("Failed to lock tokens in read pool: %v", err)
		}

		fmt.Println("locked")
	},
}

func init() {
	rootCmd.AddCommand(rpLock)
	rpLock.PersistentFlags().Float64("tokens", 0.0,
		"lock tokens number, required")
	rpLock.PersistentFlags().Float64("fee", 0.0,
		"transaction fee, default 0")
	rpLock.MarkFlagRequired("tokens")
}
