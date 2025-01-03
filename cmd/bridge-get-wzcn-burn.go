package cmd

import (
	"fmt"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
)

func init() {
	rootCmd.AddCommand(
		bridge.CreateCommandWithBridge(
			"bridge-get-wzcn-burn",
			"get confirmed burn ticket for ethereum burn transaction",
			"get transaction ticket with the given Ethereum transaction hash",
			commandGetETHBurnTicket,
			false,
			bridge.WithHash("Ethereum transaction hash"),
		))
}

func commandGetETHBurnTicket(b *zcnbridge.BridgeClient, args ...*bridge.Arg) {
	hash := bridge.GetHash(args)

	payload, err := b.QueryZChainMintPayload(hash)
	if err != nil {
		util.ExitWithError(err)
	}

	fmt.Println("WZCN burn ticket the completed consensus")
	fmt.Printf("Transaction nonce: %d\n", payload.Nonce)
	fmt.Printf("Transaction amount: %d\n", payload.Amount)
	fmt.Printf("Ethereum transaction ID: %s\n", payload.EthereumTxnID)
	fmt.Printf("ZCN receiving client ID: %s\n", payload.ReceivingClientID)
}
