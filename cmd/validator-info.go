package cmd

import (
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var validatorInfoCmd = &cobra.Command{
	Use:   "validator-info",
	Short: "Get validator info",
	Long:  `Get validator info`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags = cmd.Flags()

			json        bool
			validatorID string
			err         error
		)

		if flags.Changed("json") {
			if json, err = flags.GetBool("json"); err != nil {
				log.Fatal("invalid 'json' flag: ", err)
			}
		}

		if !flags.Changed("validator_id") {
			log.Fatal("missing required 'validator_id' flag")
		}

		if validatorID, err = flags.GetString("validator_id"); err != nil {
			log.Fatal("error in 'validator_id' flag: ", err)
		}

		var validator *sdk.Validator
		if validator, err = sdk.GetValidator(validatorID); err != nil {
			log.Fatal(err)
		}

		if json {
			util.PrintJSON(validator)
		} else {
			util.PrintValidators([]*sdk.Validator{validator})
		}

	},
}

func init() {
	rootCmd.AddCommand(validatorInfoCmd)

	validatorInfoCmd.Flags().String("validator_id", "", "validator ID, required")
	validatorInfoCmd.Flags().Bool("json", false,
		"(default false) pass this option to print response as json data")
	validatorInfoCmd.MarkFlagRequired("validator_id")
}
