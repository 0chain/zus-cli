package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var minerscInfo = &cobra.Command{
	Use:   "mn-info",
	Short: "Get miner/sharder info from Miner SC.",
	Long:  "Get miner/sharder info from Miner SC.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags = cmd.Flags()
			id    string
			err   error
			res   []byte
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if res, err = zcncore.GetMinerSCNodeInfo(id); err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(res))
	},
}

func init() {
	rootCmd.AddCommand(minerscInfo)
	minerscInfo.PersistentFlags().String("id", "", "miner/sharder ID to get info for")
	minerscInfo.MarkFlagRequired("id")
}
