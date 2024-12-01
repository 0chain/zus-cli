package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var minerscUserInfo = &cobra.Command{
	Use:   "mn-user-info",
	Short: "Get miner/sharder user pools info from Miner SC.",
	Long:  "Get miner/sharder user pools info from Miner SC.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags    = cmd.Flags()
			clientID string

			err error
		)

		if flags.Changed("client_id") {
			if clientID, err = flags.GetString("client_id"); err != nil {
				log.Fatal(err)
			}
		}

		var (
			info = new(zcncore.MinerSCUserPoolsInfo)
			res  []byte
		)
		if res, err = zcncore.GetMinerSCUserInfo(clientID); err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(res, info); err != nil {
			log.Fatal(err)
		}

		if flags.Changed("json") {
			var j bool
			if j, err = flags.GetBool("json"); err != nil {
				log.Fatal(err)
			}
			if j {
				util.PrintJSON(info)
				return
			}
		}

		if len(info.Pools) == 0 {
			fmt.Println("no user pools in Miner SC")
			return
		}

		var total common.Balance
		for _, delegates := range info.Pools {
			for _, pool := range delegates {
				total += pool.Balance
			}
		}

		for key, delegates := range info.Pools {
			for _, pool := range delegates {
				fmt.Println("- delegates:", "("+key+")")
				fmt.Println("  - pool_id:            ", pool.ID)
				fmt.Println("    balance:            ", pool.Balance)
				fmt.Println("    rewards uncollected:", pool.Reward)
				fmt.Println("    rewards paid:       ", pool.RewardPaid)
				fmt.Println("    status:             ", strings.ToLower(pool.Status))
				fmt.Println("    stake %:            ",
					float64(pool.Balance)/float64(total)*100.0, "%")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(minerscUserInfo)
	minerscUserInfo.PersistentFlags().String("client_id", "", "get info for user, if empty, current user used, optional")
	minerscUserInfo.PersistentFlags().Bool("json", false, "as JSON")
}
