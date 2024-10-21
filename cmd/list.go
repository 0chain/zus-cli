package cmd

import (
	"os"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

// listCmd represents list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list files from blobbers",
	Long:  `list files from blobbers`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags() // fflags is a *flag.FlagSet
		if fflags.Changed("remotepath") == false && fflags.Changed("authticket") == false {
			util.PrintError("Error: remotepath / authticket flag is missing")
			os.Exit(1)
		}

		remotepath := cmd.Flag("remotepath").Value.String()
		authticket := cmd.Flag("authticket").Value.String()
		lookuphash := cmd.Flag("lookuphash").Value.String()
		doJSON, _ := cmd.Flags().GetBool("json")
		if len(remotepath) == 0 && (len(authticket) == 0) {
			util.PrintError("Error: remotepath / authticket / lookuphash flag is missing")
			os.Exit(1)
		}

		if len(remotepath) > 0 {
			if fflags.Changed("allocation") == false { // check if the flag "path" is set
				util.PrintError("Error: allocation flag is missing") // If not, we'll let the user know
				os.Exit(1)                                           // and os.Exit(1)
			}

			allocationID := cmd.Flag("allocation").Value.String()
			allocationObj, err := sdk.GetAllocation(allocationID)
			if err != nil {
				util.PrintError("Error fetching the allocation", err)
				os.Exit(1)
			}

			remotepath := cmd.Flag("remotepath").Value.String()
			ref, err := allocationObj.ListDir(remotepath)
			if err != nil {
				util.PrintError(err.Error())
				os.Exit(1)
			}

			util.PrintListDirResult(doJSON, ref)
		} else if len(authticket) > 0 {
			allocationObj, err := sdk.GetAllocationFromAuthTicket(authticket)
			if err != nil {
				util.PrintError("Error fetching the allocation", err)
				os.Exit(1)
			}

			at := sdk.InitAuthTicket(authticket)
			lookuphash, err = at.GetLookupHash()
			if err != nil {
				util.PrintError("Error getting the lookuphash from authticket", err)
				os.Exit(1)
			}

			ref, err := allocationObj.ListDirFromAuthTicket(authticket, lookuphash)
			if err != nil {
				util.PrintError(err.Error())
				os.Exit(1)
			}

			util.PrintListDirResult(doJSON, ref)
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().String("allocation", "", "Allocation ID")
	listCmd.PersistentFlags().String("remotepath", "", "Remote path to list from")
	listCmd.PersistentFlags().String("authticket", "", "Auth ticket fot the file to download if you dont own it")
	listCmd.PersistentFlags().String("lookuphash", "", "The remote lookuphash of the object retrieved from the list")
	listCmd.Flags().Bool("json", false, "(default false) pass this option to print response as json data")
	listCmd.MarkFlagRequired("allocation")
}
