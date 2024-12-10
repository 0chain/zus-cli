package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zboxcore/blockchain"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/spf13/cobra"
)

var blobberUpdateCmd = &cobra.Command{
	Use:   "bl-update",
	Short: "Update blobber settings by its delegate_wallet owner",
	Long:  `Update blobber settings by its delegate_wallet owner`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()

			blobberID string
			err       error
		)

		if !flags.Changed("blobber_id") {
			log.Fatal("missing required 'blobber_id' flag")
		}

		if blobberID, err = flags.GetString("blobber_id"); err != nil {
			log.Fatal("error in 'blobber_id' flag: ", err)
		}

		if _, err = sdk.GetBlobber(blobberID); err != nil {
			log.Fatal(err)
		}

		updateBlobber := new(sdk.UpdateBlobber)
		updateBlobber.ID = common.Key(blobberID)
		if flags.Changed("capacity") {
			var capacity int64
			if capacity, err = flags.GetInt64("capacity"); err != nil {
				log.Fatal(err)
			}

			changedCapacity := common.Size(capacity)
			updateBlobber.Capacity = &changedCapacity
		}

		var delegateWallet string
		if flags.Changed("delegate_wallet") {
			if delegateWallet, err = flags.GetString("delegate_wallet"); err != nil {
				log.Fatal(err)
			}
			updateBlobber.DelegateWallet = &delegateWallet
		}

		var storageVersion int
		if flags.Changed("storage_version") {
			if storageVersion, err = flags.GetInt("storage_version"); err != nil {
				log.Fatal(err)
			}
			updateBlobber.StorageVersion = &storageVersion
		}

		terms := &sdk.UpdateTerms{}
		var termsChanged bool
		if flags.Changed("read_price") {
			var rp float64
			if rp, err = flags.GetFloat64("read_price"); err != nil {
				log.Fatal(err)
			}
			readPriceBalance, err := common.ToBalance(rp)
			if err != nil {
				log.Fatal(err)
			}
			terms.ReadPrice = &readPriceBalance
			termsChanged = true
		}

		if flags.Changed("write_price") {
			var wp float64
			if wp, err = flags.GetFloat64("write_price"); err != nil {
				log.Fatal(err)
			}
			writePriceBalance, err := common.ToBalance(wp)
			if err != nil {
				log.Fatal(err)
			}
			terms.WritePrice = &writePriceBalance
			termsChanged = true
		}

		if flags.Changed("max_offer_duration") {
			var mod time.Duration
			if mod, err = flags.GetDuration("max_offer_duration"); err != nil {
				log.Fatal(err)
			}
			terms.MaxOfferDuration = &mod
		}

		stakePoolSettings := &blockchain.UpdateStakePoolSettings{}
		var stakePoolSettingChanged bool
		if flags.Changed("num_delegates") {
			var nd int
			if nd, err = flags.GetInt("num_delegates"); err != nil {
				log.Fatal(err)
			}
			stakePoolSettings.NumDelegates = &nd
			stakePoolSettingChanged = true
		}

		if flags.Changed("delegate_Wallet") {
			var dw string
			if dw, err = flags.GetString("delegate_wallet"); err != nil {
				log.Fatal(err)
			}
			stakePoolSettings.DelegateWallet = &dw
			stakePoolSettingChanged = true
		}

		if flags.Changed("service_charge") {
			var sc float64
			if sc, err = flags.GetFloat64("service_charge"); err != nil {
				log.Fatal(err)
			}
			stakePoolSettings.ServiceCharge = &sc
			stakePoolSettingChanged = true
		}

		if flags.Changed("url") {
			var url string
			if url, err = flags.GetString("url"); err != nil {
				log.Fatal(err)
			}
			updateBlobber.BaseURL = &url
		}

		if flags.Changed("not_available") {
			var na bool
			if na, err = flags.GetBool("not_available"); err != nil {
				log.Fatal(err)
			}
			if !na {
				na = false
			}
			updateBlobber.NotAvailable = &na
		}

		if flags.Changed("is_restricted") {
			var ia bool
			// Check if the flag is set to true
			if ia, err = flags.GetBool("is_restricted"); err != nil {
				log.Fatal(err)
			}
			if !ia {
				ia = false
			}
			updateBlobber.IsRestricted = &ia
		}

		if termsChanged {
			updateBlobber.Terms = terms
		}

		if stakePoolSettingChanged {
			updateBlobber.StakePoolSettings = stakePoolSettings
		}

		if _, _, err = sdk.UpdateBlobberSettings(updateBlobber); err != nil {
			log.Fatal(err)
		}
		fmt.Println("blobber settings updated successfully")
	},
}

func init() {
	rootCmd.AddCommand(blobberUpdateCmd)
	buf := blobberUpdateCmd.Flags()
	buf.String("blobber_id", "", "blobber ID, required")
	buf.String("delegate_wallet", "", "delegate wallet, optional")
	buf.Int("storage_version", 0, "update storage version, optional")
	buf.Int64("capacity", 0, "update blobber capacity bid, optional")
	buf.Float64("read_price", 0.0, "update read_price, optional")
	buf.Float64("write_price", 0.0, "update write_price, optional")
	buf.Duration("max_offer_duration", 0*time.Second, "update max_offer_duration, optional")
	buf.Float64("min_stake", 0.0, "update min_stake, optional")
	buf.Float64("max_stake", 0.0, "update max_stake, optional")
	buf.Int("num_delegates", 0, "update num_delegates, optional")
	buf.Float64("service_charge", 0.0, "update service_charge, optional")
	buf.Bool("not_available", true, "(default false) set blobber's availability for new allocations")
	buf.Bool("is_restricted", true, "(default false) set is_restricted")
	buf.String("url", "", "update the url of the blobber, optional")
	blobberUpdateCmd.MarkFlagRequired("blobber_id")
}
