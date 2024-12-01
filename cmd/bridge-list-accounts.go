package cmd

import (
	"fmt"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
	"github.com/spf13/cobra"
)

var listEthAccounts = &cobra.Command{
	Use:   "bridge-list-accounts",
	Short: "List Ethereum account registered in local key chain",
	Long:  `List available ethereum accounts`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {
		fflags := cmd.Flags()
		path, err := fflags.GetString(bridge.OptionConfigFolder)
		if err != nil {
			fmt.Printf("Flag '%s' not found, defaulting to %s\n", bridge.OptionConfigFolder, cfgDir)
		}

		accounts := zcnbridge.ListStorageAccounts(path)
		if len(accounts) == 0 {
			fmt.Println("Accounts not found")
		}

		fmt.Println("Ethereum available account:")
		for _, acc := range accounts {
			fmt.Println(acc.Hex())
		}
	},
}

var cfgDir string

func init() {
	f := listEthAccounts
	rootCmd.AddCommand(listEthAccounts)

	var err error
	cfgDir, err = GetConfigDir()
	if err != nil {
		util.ExitWithError(err)
	}

	f.PersistentFlags().String(bridge.OptionConfigFolder, cfgDir, "Configuration dir")
}
