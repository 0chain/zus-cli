package cmd

import (
	"os"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

// The getUploadCostCmd returns value in tokens to upload a file.
var getDiffCmd = &cobra.Command{
	Use:    "get-diff",
	Short:  "Get difference of local and allocation root",
	Long:   `Get difference of local and allocation root`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {

		fflags := cmd.Flags() // fflags is a *flag.FlagSet
		if fflags.Changed("localpath") == false {
			util.PrintError("Error: localpath flag is missing")
			os.Exit(1)
		}

		localpath := cmd.Flag("localpath").Value.String()

		if len(localpath) == 0 {
			util.PrintError("Error: localpath flag is missing")
			os.Exit(1)
		}

		if fflags.Changed("allocation") == false { // check if the flag "path" is set
			util.PrintError("Error: allocation flag is missing") // If not, we'll let the user know
			os.Exit(1)                                           // and os.Exit(1)
		}
		allocationID := cmd.Flag("allocation").Value.String()

		localcache := ""
		if fflags.Changed("localcache") {
			localcache = cmd.Flag("localcache").Value.String()
		}
		exclPath := []string{}
		if fflags.Changed("excludepath") {
			exclPath, _ = cmd.Flags().GetStringArray("excludepath")
		}

		allocationObj, err := sdk.GetAllocation(allocationID)
		if err != nil {
			util.PrintError("Error fetching the allocation", err)
			os.Exit(1)
		}

		// Create filter
		filter := []string{".DS_Store", ".git"}
		lDiff, err := allocationObj.GetAllocationDiff(localcache, localpath, filter, exclPath, "/")
		if err != nil {
			util.PrintError("Error getting diff.", err)
			os.Exit(1)
		}

		util.PrintJSON(lDiff)
	},
}

func init() {
	rootCmd.AddCommand(getDiffCmd)
	getDiffCmd.PersistentFlags().String("allocation", "", "Allocation ID")
	getDiffCmd.PersistentFlags().String("localpath", "", "Local dir path to sync")
	getDiffCmd.PersistentFlags().String("localcache", "", `Local cache of remote snapshot.
If file exists, this will be used for comparison with remote.
After sync complete, remote snapshot will be updated to the same file for next use.`)
	getDiffCmd.PersistentFlags().StringArray("excludepath", []string{}, "Remote folder paths exclude to sync")
	getDiffCmd.MarkFlagRequired("allocation")
	getDiffCmd.MarkFlagRequired("localpath")
}
