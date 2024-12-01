package cmd

import (
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
)

func init() {
	command := bridge.CreateCommand(
		"bridge-verify",
		"verify ethereum transaction ",
		`verify transaction.
					<hash>`,
		VerifyEthereumTransaction,
		false,
		bridge.WithHash("Ethereum transaction hash"),
	)

	rootCmd.AddCommand(command)
}

func VerifyEthereumTransaction(args ...*bridge.Arg) {
	hash := bridge.GetHash(args)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 60, time.Second)
	if err != nil {
		util.ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("\nTransaction verification success: %s\n", hash)
	} else {
		util.ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
	}
}
