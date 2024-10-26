package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var minerscMiners = &cobra.Command{
	Use:   "ls-miners",
	Short: "Get miners from Miner SC",
	Long:  "Get miners from Miner SC",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			err   error
			info  = new(zcncore.MinerSCNodes)
		)

		limit, offset := 20, 0
		active := true
		stakable := false

		var allFlag, jsonFlag bool

		if flags.Changed("all") {
			allFlag, err = flags.GetBool("all")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("limit") {
			limit, err = flags.GetInt("limit")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("offset") {
			offset, err = flags.GetInt("offset")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("active") {
			active, err = flags.GetBool("active")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("json") {
			jsonFlag, err = flags.GetBool("json")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("stakable") {
			stakable, err = flags.GetBool("stakable")
			if err != nil {
				log.Fatal(err)
			}
		}

		if !allFlag {
			res, err := zcncore.GetMiners(active, stakable, limit, offset)
			if err != nil {
				log.Fatal(err)
			}

			if err = json.Unmarshal(res, info); err != nil {
				log.Fatal(err)
			}

			if jsonFlag {
				util.PrintJSON(info)
				return
			}

			if len(info.Nodes) == 0 {
				fmt.Println("no miners in Miner SC")
				return
			}

			util.PrintMinerNodes(info.Nodes)
			return
		} else {
			limit = 20
			offset = 0

			var nodes []zcncore.Node
			for curOff := offset; ; curOff += limit {
				res, err := zcncore.GetMiners(active, stakable, limit, offset)
				if err != nil {
					log.Fatal(err)
				}

				if err = json.Unmarshal(res, info); err != nil {
					log.Fatal(err)
				}

				if len(info.Nodes) == 0 {
					break
				}

				nodes = append(nodes, info.Nodes...)
			}

			if jsonFlag {
				util.PrintJSON(nodes)
			} else {
				util.PrintMinerNodes(nodes)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(minerscMiners)
	minerscMiners.PersistentFlags().Bool("json", false, "as JSON")
	minerscMiners.PersistentFlags().Int("limit", 20, "Limits the amount of miners returned")
	minerscMiners.PersistentFlags().Int("offset", 0, "Skips the number of miners mentioned")
	minerscMiners.PersistentFlags().Bool("active", true, "Gets active miners only, set it false to get all miners")
	minerscMiners.PersistentFlags().Bool("stakable", false, "Gets stakable miners only if set to true")
	minerscMiners.PersistentFlags().Bool("all", false, "include all registered miners, default returns the first page of miners")
}
