package cmd

import (
	"fmt"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		bridge.CreateCommandWithBridge(
			"bridge-burn-zcn",
			"burn zcn tokens",
			"burn zcn tokens that will be minted for WZCN tokens",
			commandBurnZCN,
			false,
			bridge.WithToken("ZCN tokens quantity to be burned"),
		))
}

func commandBurnZCN(b *zcnbridge.BridgeClient, args ...*bridge.Arg) {
	amount := bridge.GetToken(args)

	fmt.Println("Starting burn transaction")
	hash, _, err := b.BurnZCN(zcncore.ConvertToValue(amount))
	if err == nil {
		fmt.Printf("Submitted burn transaction %s\n", hash)
	} else {
		util.ExitWithError(err)
	}

	fmt.Printf("Transaction completed successfully: %s\n", hash)
}
