package cmd

import (
	"os"
	"strings"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var listAllCmd = &cobra.Command{
	Use:   "list-all",
	Short: "list all files from blobbers",
	Long:  `list all files from blobbers`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()                      // fflags is a *flag.FlagSet
		if fflags.Changed("allocation") == false { // check if the flag "path" is set
			util.PrintError("Error: allocation flag is missing") // If not, we'll let the user know
			os.Exit(1)                                           // and return
		}

		allocationID := cmd.Flag("allocation").Value.String()
		allocationObj, err := sdk.GetAllocation(allocationID)
		if err != nil {
			util.PrintError("Error fetching the allocation", err)
			os.Exit(1)
		}
		ref, err := allocationObj.GetRemoteFileMap(nil, "/")
		if err != nil {
			util.PrintError(err.Error())
			os.Exit(1)
		}

		type fileResp struct {
			sdk.FileInfo
			Name string `json:"name"`
			Path string `json:"path"`
		}

		fileResps := make([]fileResp, 0)
		for path, data := range ref {
			paths := strings.SplitAfter(path, "/")
			var resp = fileResp{
				Name:     paths[len(paths)-1],
				Path:     path,
				FileInfo: data,
			}
			fileResps = append(fileResps, resp)
		}

		util.PrintJSON(fileResps)
		return
	},
}

func init() {
	rootCmd.AddCommand(listAllCmd)
	listAllCmd.PersistentFlags().String("allocation", "", "Allocation ID")
	listAllCmd.MarkFlagRequired("allocation")
}
