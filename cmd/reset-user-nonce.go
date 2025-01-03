package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
)

//goland:noinspection ALL
func init() {
	rootCmd.AddCommand(
		bridge.CreateCommandWithBridge(
			"reset-user-nonce",
			"Reset user minted nonce",
			"Resets user minted nonce in bridge SC",
			resetUserNonce,
			false,
		))
}

// resetUserNonce resets user nonce in bridge SC
func resetUserNonce(bc *zcnbridge.BridgeClient, args ...*bridge.Arg) {
	tx, err := bc.ResetUserNonceMinted(context.Background())
	if err != nil {
		util.ExitWithError(err)
	}

	hash := tx.Hash().String()
	fmt.Printf("Confirming Reset of user nonce transaction: %s\n", hash)

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
