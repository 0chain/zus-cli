package cmd

import (
	"fmt"

	"github.com/0chain/gosdk/core/client"
	"github.com/0chain/gosdk/core/encryption"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign-data",
	Short: "Sign given data",
	Long:  `Sign given data`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		data, _ := cmd.Flags().GetString("data")
		if data == "" {
			data = client.Id()
		} else {
			data = encryption.Hash(data)
		}
		sign, err := client.Sign(data)
		if err != nil {
			fmt.Println("Error generating the signature. ", err.Error())
			return
		}
		fmt.Println("Signature : " + sign)
		return
	},
}

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String("data", "", "give data for signing, Default will be clientID")
}
