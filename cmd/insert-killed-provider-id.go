package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/spf13/cobra"
)

var insertKilledProviderId = &cobra.Command{
	Use:   "insert-killed-provider-id",
	Short: "Insert killed provider id",
	Long:  `Insert killed provider id`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()

			blobberID string
			err       error
		)

		if !flags.Changed("id") {
			log.Fatal("missing required 'blobber_id' flag")
		}
		if blobberID, err = flags.GetString("id"); err != nil {
			log.Fatal("error in 'id' flag: ", err)
		}

		snv := sdk.StorageNodeIdField{
			Id: blobberID,
		}

		_, _, err = sdk.InsertKilledProviderID(&snv)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("insert killed id successfully")
	},
}

func init() {
	rootCmd.AddCommand(insertKilledProviderId)
	insertKilledProviderId.Flags().String("id", "", "blobber_id is required")
	insertKilledProviderId.MarkFlagRequired("id")
}
