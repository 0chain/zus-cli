package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var createWalletCmd = &cobra.Command{
	Use:   "create-wallet",
	Short: "Create wallet and logs it into stdout (pass --register to register wallet to blockchain)",
	Long:  `Create wallet and logs it into standard output (pass --register to register wallet to blockchain)`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		walletStr, err := zcncore.CreateWalletOffline()
		if err != nil {
			util.ExitWithError("failed to generate offline wallet", err)
		}
		walletName := cmd.Flags().Lookup("wallet").Value.String()

		// write wallet into wallet dir
		filename := walletFilename(walletName)
		if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
			// wallet exists
			fmt.Printf("wallet already exists at %s", filename)
			return
		}

		if err := os.WriteFile(filename, []byte(walletStr), 0644); err != nil {
			// no return just print it
			fmt.Fprintf(os.Stderr, "failed to dump wallet into zcn home directory %v", err)
		} else {
			log.Printf("wallet saved in %s\n", filename)
		}

		if !bSilent {
			fmt.Fprint(os.Stdout, walletStr)
		}
	},
}

func init() {
	rootCmd.AddCommand(WithoutWallet(createWalletCmd))
	createWalletCmd.PersistentFlags().Bool("silent", false, "do not print wallet details in the standard output (default false)")
	createWalletCmd.PersistentFlags().String("wallet", "", "give custom name to the wallet")
}

func walletFilename(walletName string) string {
	cfgDir, err := GetConfigDir()
	if err != nil {
		util.ExitWithError(err)
	}
	if len(walletName) > 0 {
		return filepath.Join(cfgDir, walletName)
	}
	return filepath.Join(cfgDir, "wallet.json")
}
