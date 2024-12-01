package cmd

import (
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
	"github.com/spf13/cobra"
)

//goland:noinspection GoUnhandledErrorResult
func init() {
	cmd := &cobra.Command{
		Use:   "bridge-import-account",
		Short: "Import Ethereum account to local key storage (default $HOME/.zcn/wallets)",
		Long:  "Import account to local key storage using mnemonic, protected with password (default $HOME/.zcn/wallets)",
		Args:  cobra.MinimumNArgs(0),
		Run: func(c *cobra.Command, _ []string) {
			bridge.Check(c, bridge.OptionMnemonic, bridge.OptionKeyPassword)

			path := c.Flag(bridge.OptionConfigFolder).Value.String()
			mnemonic := c.Flag(bridge.OptionMnemonic).Value.String()
			password := c.Flag(bridge.OptionKeyPassword).Value.String()
			var accountAddrIndex zcnbridge.AccountAddressIndex

			if c.Flags().Changed(bridge.OptionAccountIndex) {
				var err error
				accountAddrIndex.AccountIndex, err = c.Flags().GetInt(bridge.OptionAccountIndex)
				if err != nil {
					util.ExitWithError(err)
					return
				}
			}

			if c.Flags().Changed(bridge.OptionAddressIndex) {
				var err error
				accountAddrIndex.AddressIndex, err = c.Flags().GetInt(bridge.OptionAddressIndex)
				if err != nil {
					util.ExitWithError(err)
					return
				}
			}

			if c.Flags().Changed(bridge.OptionBip32) {
				var err error
				accountAddrIndex.Bip32, err = c.Flags().GetBool(bridge.OptionBip32)
				if err != nil {
					util.ExitWithError(err)
					return
				}
			}

			_, err := zcnbridge.ImportAccount(path, mnemonic, password, accountAddrIndex)
			if err != nil {
				util.ExitWithError(err)
				return
			}
		},
	}

	cfgDir, err := GetConfigDir()
	if err != nil {
		util.ExitWithError(err)
	}

	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String(bridge.OptionMnemonic, "", "Ethereum mnemonic")
	cmd.PersistentFlags().String(bridge.OptionKeyPassword, "", "Password to lock and unlock account to sign transaction")
	cmd.PersistentFlags().Int(bridge.OptionAccountIndex, 0, "Index of the account to use, default 0")
	cmd.PersistentFlags().Int(bridge.OptionAddressIndex, 0, "Index of the address to use, default 0")
	cmd.PersistentFlags().Bool(bridge.OptionBip32, false, "Use BIP32 derivation path")
	cmd.PersistentFlags().String(bridge.OptionConfigFolder, cfgDir, "Home config directory")

	cmd.MarkFlagRequired(bridge.OptionMnemonic)
	cmd.MarkFlagRequired(bridge.OptionKeyPassword)
}
