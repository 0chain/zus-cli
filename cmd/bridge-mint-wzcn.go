package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
)

func init() {
	rootCmd.AddCommand(
		bridge.CreateCommandWithBridge(
			"bridge-mint-wzcn",
			"mint WZCN tokens using the hash of ZCN burn transaction",
			"mint WZCN tokens after burning ZCN tokens in ZCN chain",
			commandMintEth,
			false,
		))
}

func commandMintEth(b *zcnbridge.BridgeClient, args ...*bridge.Arg) {
	userNonce, err := b.GetUserNonceMinted(context.Background(), b.EthereumAddress)
	if err != nil {
		util.ExitWithError(err)
	}

	var burnTickets []zcncore.BurnTicket

	res, err := zcncore.GetNotProcessedZCNBurnTickets(b.EthereumAddress, userNonce.String())
	if err != nil {
		util.ExitWithError(err)
	}

	err = json.Unmarshal(res, &burnTickets)
	if err != nil {
		util.ExitWithError(err)
	}

	fmt.Printf("Found %d not processed ZCN burn transactions\n", len(burnTickets))

	for _, burnTicket := range burnTickets {
		fmt.Printf("Query ticket for ZCN transaction hash: %s\n", burnTicket.Hash)

		payload, err := b.QueryEthereumMintPayload(burnTicket.Hash)
		if err != nil {
			util.ExitWithError(err)
		}

		fmt.Printf("Sending mint transaction to Ethereum\n")
		fmt.Printf("Payload amount: %d\n", payload.Amount)
		fmt.Printf("Payload nonce: %d\n", payload.Nonce)
		fmt.Printf("ZCN transaction ID: %s\n", payload.ZCNTxnID)

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*20)
		defer cancelFunc()

		fmt.Println("Starting to mint WZCN")

		tx, err := b.MintWZCN(ctx, payload)
		if err != nil {
			util.ExitWithError(err)
		}

		hash := tx.Hash().String()
		fmt.Printf("Confirming Ethereum mint transaction: %s\n", hash)

		status, err := zcnbridge.ConfirmEthereumTransaction(hash, 20, time.Second*5)
		if err != nil {
			util.ExitWithError(err)
		}

		if status == 1 {
			fmt.Printf("\nTransaction verification success: %s\n", hash)
		} else {
			util.ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
		}
	}

	if len(burnTickets) > 0 {
		fmt.Println("Done.")
	} else {
		fmt.Println("Failed.")
	}
}
