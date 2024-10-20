package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/spf13/cobra"
)

var resetBlobberStatsCmd = &cobra.Command{
	Use:    "reset-blobber-stats",
	Short:  "Reset blobber stats",
	Long:   `Reset blobber stats`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()

			blobberID     string
			prevAllocated int64
			prevSavedData int64
			newAllocated  int64
			newSavedData  int64
			err           error
		)

		if !flags.Changed("blobber_id") {
			log.Fatal("missing required 'blobber_id' flag")
		}
		if blobberID, err = flags.GetString("blobber_id"); err != nil {
			log.Fatal("error in 'blobber_id' flag: ", err)
		}

		if !flags.Changed("prev_allocated") {
			log.Fatal("missing required 'prev_allocated' flag")
		}
		if prevAllocated, err = flags.GetInt64("prev_allocated"); err != nil {
			log.Fatal("error in 'prev_allocated' flag: ", err)
		}

		if !flags.Changed("prev_saved_data") {
			log.Fatal("missing required 'prev_saved_data' flag")
		}
		if prevSavedData, err = flags.GetInt64("prev_saved_data"); err != nil {
			log.Fatal("error in 'prev_saved_data' flag: ", err)
		}

		if !flags.Changed("new_allocated") {
			log.Fatal("missing required 'new_allocated' flag")
		}
		if newAllocated, err = flags.GetInt64("new_allocated"); err != nil {
			log.Fatal("error in 'new_allocated' flag: ", err)
		}

		if !flags.Changed("new_saved_data") {
			log.Fatal("missing required 'new_saved_data' flag")
		}
		if newSavedData, err = flags.GetInt64("new_saved_data"); err != nil {
			log.Fatal("error in 'new_saved_data' flag: ", err)
		}

		resetBlobberStatsDto := &sdk.ResetBlobberStatsDto{
			BlobberID:     blobberID,
			PrevAllocated: prevAllocated,
			PrevSavedData: prevSavedData,
			NewAllocated:  newAllocated,
			NewSavedData:  newSavedData,
		}
		fmt.Println(*resetBlobberStatsDto)

		_, _, err = sdk.ResetBlobberStats(resetBlobberStatsDto)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("reset blobber stats successfully")
	},
}

func init() {
	rootCmd.AddCommand(resetBlobberStatsCmd)
	resetBlobberStatsCmd.Flags().String("blobber_id", "", "blobber_id is required")
	resetBlobberStatsCmd.Flags().Int64("prev_allocated", 0, "prev_allocated is required")
	resetBlobberStatsCmd.Flags().Int64("prev_saved_data", 0, "prev_saved_data is required")
	resetBlobberStatsCmd.Flags().Int64("new_allocated", 0, "new_allocated is required")
	resetBlobberStatsCmd.Flags().Int64("new_saved_data", 0, "new_saved_data is required")
	resetBlobberStatsCmd.MarkFlagRequired("blobber_id")
	resetBlobberStatsCmd.MarkFlagRequired("prev_allocated")
	resetBlobberStatsCmd.MarkFlagRequired("prev_saved_data")
	resetBlobberStatsCmd.MarkFlagRequired("new_allocated")
	resetBlobberStatsCmd.MarkFlagRequired("new_saved_data")
}
