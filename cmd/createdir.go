package cmd

import (
	"fmt"
	"os"

	"github.com/0chain/gosdk/constants"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var createDirCmd = &cobra.Command{
	Use:   "createdir",
	Short: "Create directory",
	Long:  `Create directory`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()              // fflags is a *flag.FlagSet
		if !fflags.Changed("allocation") { // check if the flag "path" is set
			util.PrintError("Error: allocation flag is missing") // If not, we'll let the user know
			os.Exit(1)                                           // and return
		}
		if !fflags.Changed("dirname") {
			util.PrintError("Error: dirname flag is missing")
			os.Exit(1)
		}

		allocationID := cmd.Flag("allocation").Value.String()
		allocationObj, err := sdk.GetAllocation(allocationID)
		if err != nil {
			util.PrintError("Error fetching the allocation", err)
			os.Exit(1)
		}
		dirname := cmd.Flag("dirname").Value.String()

		err = allocationObj.DoMultiOperation([]sdk.OperationRequest{
			{
				OperationType: constants.FileOperationCreateDir,
				RemotePath:    dirname,
			},
		})

		if err != nil {
			util.PrintError("Directory creation failed.", err)
			os.Exit(1)
		}

		fmt.Println(dirname + " directory created")
	},
}

func init() {

	rootCmd.AddCommand(createDirCmd)
	createDirCmd.PersistentFlags().String("allocation", "", "Allocation ID")
	createDirCmd.PersistentFlags().String("dirname", "", "New directory name")
	createDirCmd.MarkFlagRequired("allocation")
	createDirCmd.MarkFlagRequired("dirname")
}
