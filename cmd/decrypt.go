package cmd

import (
	"fmt"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var walletDecryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt text with passphrase",
	Long:  `Decrypt text with passphrase`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		passphrase, _ := cmd.Flags().GetString("passphrase")
		text, _ := cmd.Flags().GetString("text")

		decrypted, err := zcncore.Decrypt(passphrase, text)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(decrypted)
		return
	},
}

func init() {
	rootCmd.AddCommand(walletDecryptCmd)
	walletDecryptCmd.Flags().String("passphrase", "", "Passphrase to decrypt text")
	walletDecryptCmd.Flags().String("text", "", "Encrypted text")
}
