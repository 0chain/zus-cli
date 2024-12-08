package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/spf13/cobra"
)

var resetVersionCmd = &cobra.Command{
	Use:   "reset-version",
	Short: "Reset blobber version",
	Long:  `Reset blobber version`,
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

		snv := sdk.StorageNodeIdField{
			Id: blobberID,
		}

		_, _, err = sdk.ResetBlobberVersion(&snv)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("reset blobber version successfully")
	},
}

func init() {
	rootCmd.AddCommand(resetVersionCmd)
	resetVersionCmd.Flags().String("blobber_id", "", "blobber_id is required")
	resetVersionCmd.MarkFlagRequired("blobber_id")

}
