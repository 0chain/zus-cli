package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
	comm "github.com/ethereum/go-ethereum/common"
)

//goland:noinspection ALL
func init() {
	rootCmd.AddCommand(
		bridge.CreateCommandWithBridge(
			"auth-sc-delete",
			"Deletes an authorizer to token bridge SC manually",
			"Deletes an authorizer to token bridge SC manually",
			deleteAuthorizerInSC,
			true,
			&bridge.Option{
				Name:     "ethereum_address",
				Typename: "string",
				Value:    "",
				Usage:    "ethereum address which is authorizer linked to",
				Required: true,
			},
		))
}

// registerAuthorizerInSC registers a new authorizer to token bridge SC
func deleteAuthorizerInSC(bc *zcnbridge.BridgeClient, args ...*bridge.Arg) {
	ethereumAddress := bridge.GetEthereumAddress(args)

	tx, err := bc.RemoveEthereumAuthorizer(context.Background(), comm.HexToAddress(ethereumAddress))
	if err != nil {
		util.ExitWithError(err)
	}

	hash := tx.Hash().String()
	fmt.Printf("Confirming Ethereum mint transaction: %s\n", hash)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 100, time.Second*5)
	if err != nil {
		util.ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("\nTransaction verification success: %s\n", hash)
	} else {
		util.ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
	}
}
