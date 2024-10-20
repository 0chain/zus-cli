package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var blobberInfoCmd = &cobra.Command{
	Use:   "bl-info",
	Short: "Get blobber info",
	Long:  `Get blobber info`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags = cmd.Flags()

			json      bool
			blobberID string
			err       error
		)

		if flags.Changed("json") {
			if json, err = flags.GetBool("json"); err != nil {
				log.Fatal("invalid 'json' flag: ", err)
			}
		}

		if !flags.Changed("blobber_id") {
			log.Fatal("missing required 'blobber_id' flag")
		}

		if blobberID, err = flags.GetString("blobber_id"); err != nil {
			log.Fatal("error in 'blobber_id' flag: ", err)
		}

		var blob *sdk.Blobber
		if blob, err = sdk.GetBlobber(blobberID); err != nil {
			log.Fatal(err)
		}

		if json {
			util.PrintJSON(blob)
			return
		}

		fmt.Println("id:               ", blob.ID)
		fmt.Println("url:              ", blob.BaseURL)
		fmt.Println("capacity:         ", blob.Capacity)
		fmt.Println("is killed:        ", blob.IsKilled)
		fmt.Println("is shut down:     ", blob.IsShutdown)
		fmt.Println("last_health_check:", blob.LastHealthCheck.ToTime())
		fmt.Println("capacity_used:    ", blob.Allocated)
		fmt.Println("total_stake:      ", blob.TotalStake)
		fmt.Println("not_available:     ", blob.NotAvailable)
		fmt.Println("is_restricted:     ", blob.IsRestricted)
		fmt.Println("terms:")
		fmt.Println("  read_price:        ", blob.Terms.ReadPrice, "/ GB")
		fmt.Println("  write_price:       ", blob.Terms.WritePrice, "/ GB")
		fmt.Println("  max_offer_duration:", blob.Terms.MaxOfferDuration)
		fmt.Println("settings:")
		fmt.Println("  delegate_wallet:", blob.StakePoolSettings.DelegateWallet)
		//fmt.Println("  min_stake:      ", blob.StakePoolSettings.MinStake)
		//fmt.Println("  max_stake:      ", blob.StakePoolSettings.MaxStake)
		fmt.Println("  num_delegates:  ", blob.StakePoolSettings.NumDelegates)
		fmt.Println("  service_charge: ", blob.StakePoolSettings.ServiceCharge*100, "%")
	},
}

func init() {
	rootCmd.AddCommand(blobberInfoCmd)
	blobberInfoCmd.Flags().String("blobber_id", "", "blobber ID, required")
	blobberInfoCmd.Flags().Bool("json", false,
		"(default false) pass this option to print response as json data")
	blobberInfoCmd.MarkFlagRequired("blobber_id")
}
