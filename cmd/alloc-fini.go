package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/spf13/cobra"
)

func isFinalized(allocID string) (ok bool, err error) {
	var alloc *sdk.Allocation
	if alloc, err = sdk.GetAllocation(allocID); err != nil {
		return false, fmt.Errorf("can't get allocation from sharders: %v", err)
	}
	return alloc.Finalized, nil
}

func allocShouldNotBeFinalized(allocID string) {
	var ok, err = isFinalized(allocID)
	if err != nil {
		log.Fatalf("can't get allocation from sharders: %v", err)
	}
	if ok {
		log.Fatal("allocation already finalized")
	}
}

// finiAllocationCmd used to change allocation size and expiration
var finiAllocationCmd = &cobra.Command{
	Use:   "alloc-fini",
	Short: "Finalize an expired allocation",
	Long: `Finalize an expired allocation by allocation owner or one of
blobbers of the allocation. It moves all tokens have to be moved between pools
and empties write pool moving left tokens to client.`,
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

		// check out allocation first
		allocShouldNotBeFinalized(allocID)

		txnHash, _, err := sdk.FinalizeAllocation(allocID)
		if err != nil {
			// check again, a blobber can finalize it
			allocShouldNotBeFinalized(allocID)
			// finalizing error
			log.Fatal("Error finalizing allocation:", err)
		}
		// success
		log.Print("Allocation finalized with txId : " + txnHash)
	},
}

func init() {
	rootCmd.AddCommand(finiAllocationCmd)

	finiAllocationCmd.PersistentFlags().String("allocation", "",
		"Allocation ID")
	finiAllocationCmd.MarkFlagRequired("allocation")
}
