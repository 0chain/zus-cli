package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var repairSize = &cobra.Command{
	Use:   "repair-size",
	Short: "gets only size to repair file to blobbers",
	Long:  `gets only size to repair file to blobbers`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if !fflags.Changed("allocation") {
			util.PrintError("Error: allocation flag is missing")
			os.Exit(1)
		}

		repairPath := "/"
		var err error
		if fflags.Changed("repairpath") {
			if repairPath, err = fflags.GetString("repairpath"); err != nil {
				util.PrintError("Error: repairpath is not of string type", err)
				os.Exit(1)
			}
		}

		allocationID := cmd.Flag("allocation").Value.String()
		allocationObj, err := sdk.GetAllocation(allocationID)
		if err != nil {
			util.PrintError("Error fetching the allocation.", err)
			os.Exit(1)
		}

		size, err := allocationObj.RepairSize(repairPath)
		if err != nil {
			util.PrintError("get repair size failed: ", err)
			os.Exit(1)
		}

		jsonBytes, err := json.Marshal(size)
		if err != nil {
			util.PrintError("error marshaling size: ", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonBytes))
	},
}

func init() {
	rootCmd.AddCommand(repairSize)
	repairSize.PersistentFlags().String("allocation", "", "Allocation ID")
	repairSize.PersistentFlags().String("repairpath", "", "Path to repair")
	repairSize.MarkFlagRequired("allocation")
}
