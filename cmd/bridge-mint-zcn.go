package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
)

func init() {
	rootCmd.AddCommand(
		bridge.CreateCommandWithBridge(
			"bridge-mint-zcn",
			"mint zcn tokens using the hash of Ethereum burn transaction",
			"mint zcn tokens after burning WZCN tokens in Ethereum chain",
			commandMintZCN,
			false,
			&bridge.Option{
				Name:     "burn-txn-hash",
				Typename: "string",
				Value:    "",
				Usage:    "mint the ZCN tokens for the given Ethereum burn transaction hash",
			},
		))
}

func commandMintZCN(b *zcnbridge.BridgeClient, args ...*bridge.Arg) {
	burnHash := bridge.GetString(args, "burn-txn-hash")

	var mintNonce int64
	res, err := zcncore.GetMintNonce()
	if err != nil {
		util.ExitWithError(err)
	}

	err = json.Unmarshal(res, &mintNonce)
	if err != nil {
		util.ExitWithError(err)
	}

	burnTickets, err := b.QueryEthereumBurnEvents(strconv.Itoa(int(mintNonce)))
	if err != nil {
		util.ExitWithError(err)
	}

	fmt.Printf("Found %d not processed WZCN burn transactions\n", len(burnTickets))

	for _, burnTicket := range burnTickets {
		if len(burnHash) > 0 {
			if burnHash != burnTicket.TransactionHash {
				continue
			}
		}

		fmt.Printf("Query ticket for Ethereum transaction hash: %s\n", burnTicket.TransactionHash)

		payload, err := b.QueryZChainMintPayload(burnTicket.TransactionHash)
		if err != nil {
			util.ExitWithError(err)
		}

		fmt.Printf("Sending mint transaction to ZCN\n")
		fmt.Printf("Ethereum transaction ID: %s\n", payload.EthereumTxnID)
		fmt.Printf("Payload amount: %d\n", payload.Amount)
		fmt.Printf("Payload nonce: %d\n", payload.Nonce)
		fmt.Printf("Receiving ZCN ClientID: %s\n", payload.ReceivingClientID)

		fmt.Println("Starting to mint ZCN")

		txHash, err := b.MintZCN(payload)
		if err != nil {
			util.ExitWithError(err)
		}

		fmt.Println("Completed ZCN mint transaction")
		fmt.Printf("Transaction hash: %s\n", txHash)

	}

	if len(burnTickets) > 0 {
		fmt.Println("Done.")
	} else {
		fmt.Println("Failed.")
	}
}
