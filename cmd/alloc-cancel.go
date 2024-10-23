package cmd

import (
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/spf13/cobra"
)

// cancelAllocationCmd used to cancel allocation where blobbers
// doesn't provides their service in reality
var cancelAllocationCmd = &cobra.Command{
	Use:   "alloc-cancel",
	Short: "Cancel an allocation",
	Long: `Cancel allocation used to terminate an allocation where, because
of blobbers, it can't be used. Thus, the blobbers will not receive their
min_lock_demand. Other aspects of the cancellation follows the finalize
allocation flow.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var flags = cmd.Flags()

		if flags.Changed("allocation") == false {
			log.Fatal("Error: allocation flag is missing")
		}

		allocID, err := flags.GetString("allocation")
		if err != nil {
			log.Fatal("invalid 'allocation' flag: ", err)
		}

		txnHash, _, err := sdk.CancelAllocation(allocID)
		if err != nil {
			log.Fatal("Error canceling allocation:", err)
		}
		log.Println("Allocation canceled with txId : " + txnHash)
	},
}

func init() {
	rootCmd.AddCommand(cancelAllocationCmd)
	cancelAllocationCmd.PersistentFlags().String("allocation", "",
		"Allocation ID")
	cancelAllocationCmd.MarkFlagRequired("allocation")
}
