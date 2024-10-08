package cmd

import (
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

// spUserInfo information per user
var spUserInfo = &cobra.Command{
	Use:   "sp-user-info",
	Short: "Stake pool information for a user.",
	Long:  `Stake pool information for a user.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags    = cmd.Flags()
			limit    int
			offset   int
			isAll    bool
			clientID string
			err      error
		)

		doJSON, _ := cmd.Flags().GetBool("json")

		if flags.Changed("client_id") {
			if clientID, err = flags.GetString("client_id"); err != nil {
				log.Fatalf("can't get 'client_id' flag: %v", err)
			}
		}

		limit, err = flags.GetInt("limit")
		if err != nil {
			log.Fatal(err)
		}

		offset, err = flags.GetInt("offset")
		if err != nil {
			log.Fatal(err)
		}

		if flags.Changed("all") {
			isAll, err = flags.GetBool("all")
			if err != nil {
				log.Fatal(err)
			}
		}

		if !isAll {
			if _, err := getAndPrintStakePool(clientID, doJSON, offset, limit); err != nil {
				log.Fatalf("Failed to get stake pool info: %v", err)
			}
			return
		}

		for curOff := offset; ; curOff += limit {
			l, err := getAndPrintStakePool(clientID, doJSON, curOff, limit)
			if err != nil {
				log.Fatalf("Failed to get stake pool info: %v", err)
			}
			if l == 0 {
				return
			}
		}

	},
}

func getAndPrintStakePool(clientID string, doJSON bool, offset, limit int) (int, error) {
	var info *sdk.StakePoolUserInfo
	var err error
	if info, err = sdk.GetStakePoolUserInfo(clientID, offset, limit); err != nil {
		return 0, err
	}
	if doJSON {
		util.PrintJSON(info)
	} else {
		util.PrintStakePoolUserInfo(info)
	}
	return len(info.Pools), nil
}

func init() {
	rootCmd.AddCommand(spUserInfo)
	spUserInfo.PersistentFlags().Bool("json", false, "(default false) pass this option to print response as json data")
	spUserInfo.PersistentFlags().Bool("all", false, "(default false) pass this option to get all the pools")
	spUserInfo.PersistentFlags().Int("limit", 20, "pass this option to limit the number of records returned")
	spUserInfo.PersistentFlags().Int("offset", 0, "pass this option to skip the number of rows before beginning")
	spUserInfo.PersistentFlags().String("client_id", "", "pass for given client")
}
