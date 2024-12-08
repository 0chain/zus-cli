package cmd

import (
	"fmt"
	"os"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var recoverwalletcmd = &cobra.Command{
	Use:   "recoverwallet",
	Short: "Recover wallet",
	Long:  `Recover wallet from the previously stored mnemonic`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("mnemonic") == false {
			util.ExitWithError("Error: Mnemonic not provided")
		}
		mnemonic := cmd.Flag("mnemonic").Value.String()
		if zcncore.IsMnemonicValid(mnemonic) == false {
			util.ExitWithError("Error: Invalid mnemonic")
		}

		walletString, err := zcncore.RecoverOfflineWallet(mnemonic)
		if err != nil {
			util.ExitWithError(err.Error())
		}

		var walletFilePath string
		cfgDir, err := GetConfigDir()
		if err != nil {
			util.ExitWithError(err.Error())
		}
		if len(walletFile) > 0 {
			walletFilePath = cfgDir + "/" + walletFile
		} else {
			walletFilePath = cfgDir + "/wallet.json"
		}
		clientConfig := string(walletString)
		file, err := os.Create(walletFilePath)
		if err != nil {
			util.ExitWithError(err.Error())
		}
		defer file.Close()
		fmt.Fprintf(file, clientConfig)
		fmt.Println("Wallet recovered!!")
		return
	},
}

func init() {
	rootCmd.AddCommand(WithoutWallet(recoverwalletcmd))
	recoverwalletcmd.PersistentFlags().String("mnemonic", "", "mnemonic")
	recoverwalletcmd.PersistentFlags().Bool("offline", false, "recover wallet without registration on blockchain")
	recoverwalletcmd.MarkFlagRequired("mnemonic")
}
