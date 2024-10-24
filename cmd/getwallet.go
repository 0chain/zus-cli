package cmd

import (
	"fmt"
	"os"

	"github.com/0chain/gosdk/core/client"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

// walletinfo used for getting the wallet info
var walletinfoCmd = &cobra.Command{
	Use:   "getwallet",
	Short: "Get wallet information",
	Long:  `Get wallet information`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		doJSON, _ := cmd.Flags().GetBool("json")

		header := []string{"Public Key", "ClientID", "Encryption Public Key"}
		data := make([][]string, 1)
		encPubKey, err := sdk.GetClientEncryptedPublicKey()
		if err != nil {
			fmt.Println("Error getting the public key for encryption. ", err.Error())
			return
		}
		data[0] = []string{client.PublicKey(), client.Id(), encPubKey}
		if doJSON {
			j := make(map[string]string)
			j["client_public_key"] = client.PublicKey()
			j["client_id"] = client.Id()
			j["encryption_public_key"] = encPubKey
			util.PrintJSON(j)
			return
		}
		util.WriteTable(os.Stdout, header, []string{}, data)
		return
	},
}

func init() {
	rootCmd.AddCommand(walletinfoCmd)
	walletinfoCmd.Flags().Bool("json", false, "(default false) pass this option to print response as json data")
}
