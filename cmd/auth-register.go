package cmd

import (
	"log"
	"strings"

	"github.com/0chain/gosdk/core/transaction"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/0chain/zus-cli/util/bridge"
	"github.com/pkg/errors"
)

//goland:noinspection ALL
func init() {
	rootCmd.AddCommand(
		bridge.CreateCommandWithBridge(
			"auth-register",
			"Register an authorizer manually",
			"Register an authorizer manually",
			registerAuthorizerInChain,
			true,
			&bridge.Option{
				Name:     "url",
				Typename: "string",
				Value:    "",
				Usage:    "authorizer endpoint url",
				Required: true,
			},
			&bridge.Option{
				Name:     "client_id",
				Typename: "string",
				Value:    "",
				Usage:    "the client_id of the wallet",
				Required: true,
			},
			&bridge.Option{
				Name:     "client_key",
				Typename: "string",
				Value:    "",
				Usage:    "the client_key which is the public key of the wallet",
				Required: true,
			},
			&bridge.Option{
				Name:     "min_stake",
				Typename: "int64",
				Value:    int64(1),
				Usage:    "the minimum stake value for the stake pool",
				Required: false,
			},
			&bridge.Option{
				Name:     "max_stake",
				Typename: "int64",
				Value:    int64(10),
				Usage:    "the maximum stake value for the stake pool",
				Required: false,
			},
			&bridge.Option{
				Name:     "num_delegates",
				Typename: "int",
				Value:    5,
				Usage:    "the number of delegates in the authorizer stake pool",
				Required: false,
			},
			&bridge.Option{
				Name:     "service_charge",
				Typename: "float64",
				Value:    0.0,
				Usage:    "the service charge for the authorizer stake pool",
				Required: false,
			},
		))
}

// registerAuthorizerInChain registers a new authorizer
// addAuthorizerPayload *addAuthorizerPayload
func registerAuthorizerInChain(bc *zcnbridge.BridgeClient, args ...*bridge.Arg) {
	clientID := bridge.GetClientID(args)
	clientKey := bridge.GetClientKey(args)
	url := bridge.GetURL(args)
	numDelegates := bridge.GetNumDelegates(args)
	serviceCharge := bridge.GetServiceCharge(args)

	input := &zcncore.AddAuthorizerPayload{
		PublicKey: clientKey,
		URL:       url,
		StakePoolSettings: zcncore.AuthorizerStakePoolSettings{
			DelegateWallet: clientID,
			NumDelegates:   numDelegates,
			ServiceCharge:  serviceCharge,
		},
	}

	hash, _, _, txn, err := zcncore.ZCNSCAddAuthorizer(input)
	if err != nil {
		log.Fatal(err, "failed to add authorizer with transaction: '%s'", hash)
	}

	log.Printf("Authorizer submitted OK... " + hash)
	log.Printf("Starting verification: " + hash)

	txn, err = transaction.VerifyTransaction(hash)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			util.ExitWithError("Authorizer has already been added to 0Chain...  Continue")
		} else {
			util.ExitWithError(errors.Wrapf(err, "failed to verify transaction: '%s'", txn.Hash))
		}
	}
}
