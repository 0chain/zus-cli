package cmd

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
)

func init() {
	rootCmd.AddCommand(
		bridge.CreateCommandWithBridge(
			"bridge-burn-eth",
			"burn eth tokens",
			"burn eth tokens that will be minted on ZCN chain",
			commandBurnEth,
			false,
			bridge.WithAmount("WZCN token amount to be burned"),
			bridge.WithRetries("Num of seconds a transaction status check should run"),
		))
}

func commandBurnEth(b *zcnbridge.BridgeClient, args ...*bridge.Arg) {
	retries := bridge.GetRetries(args)
	amount := bridge.GetAmount(args)

	var (
		transaction *types.Transaction
		err         error
		hash        string
		status      int
	)

	var balanceRaw *big.Int

	balanceRaw, err = b.GetTokenBalance()
	if err != nil {
		util.ExitWithError(err, "failed to GetTokenBalance")
	}

	balance := balanceRaw.Uint64()

	if balance < amount {
		target := amount - balance

		var estimated *big.Int

		estimated, err = b.GetETHSwapAmount(context.Background(), target)
		if err != nil {
			util.ExitWithError(err, "failed to GetETHSwapAmount")
		}

		transaction, err = b.SwapETH(context.Background(), estimated.Uint64(), target)
		if err != nil {
			util.ExitWithError(err, "failed to execute SwapETH")
		}

		hash = transaction.Hash().Hex()
		status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
		if err != nil {
			util.ExitWithError(fmt.Sprintf("Failed to confirm SwapETH: hash = %s, error = %v", hash, err))
		}

		if status == 1 {
			fmt.Printf("Verification: SwapETH [OK]: %s\n", hash)
		} else {
			util.ExitWithError(fmt.Sprintf("Verification: SwapETH [FAILED]: %s\n", hash))
		}
	}

	fmt.Println("Starting IncreaseBurnerAllowance transaction")
	transaction, err = b.IncreaseBurnerAllowance(context.Background(), amount)
	if err != nil {
		util.ExitWithError(err, "failed to execute IncreaseBurnerAllowance")
	}

	hash = transaction.Hash().Hex()
	status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
	if err != nil {
		util.ExitWithError(fmt.Sprintf("Failed to confirm IncreaseBurnerAllowance: hash = %s, error = %v", hash, err))
	}

	if status == 1 {
		fmt.Printf("Verification: IncreaseBurnerAllowance [OK]: %s\n", hash)
	} else {
		util.ExitWithError(fmt.Sprintf("Verification: IncreaseBurnerAllowance [FAILED]: %s\n", hash))
	}

	fmt.Println("Starting WZCN burn transaction")

	transaction, err = b.BurnWZCN(context.Background(), amount)
	if err != nil {
		util.ExitWithError(err, "failed to burn WZCN tokens")
	}
	hash = transaction.Hash().String()
	fmt.Printf("Confirming WZCN burn transaction %s\n", hash)

	status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
	if err != nil {
		util.ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("Verification: WZCN burn [OK]: %s\n", hash)
	}

	if status == 0 {
		util.ExitWithError(fmt.Sprintf("Verification: WZCN burn [PENDING]: %s\n", hash))
	}

	if status == -1 {
		util.ExitWithError(fmt.Sprintf("Verification: WZCN burn not started, please, check later [FAILED]: %s\n", hash))
	}
}
