package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

var minerscSharders = &cobra.Command{
	Use:   "ls-sharders",
	Short: "Get sharders from Miner SC",
	Long:  "Get sharders from Miner SC",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		flags := cmd.Flags()

		var err error
		var jsonFlag, allFlag, stakable bool

		if flags.Changed("json") {
			jsonFlag, err = flags.GetBool("json")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("all") {
			allFlag, err = flags.GetBool("all")
			if err != nil {
				log.Fatal(err)
			}
		}

		mb, err := zcncore.GetLatestFinalizedMagicBlock()
		if err != nil {
			log.Fatalf("Failed to get MagicBlock: %v", err)
		}

		limit, offset := 20, 0
		active := true
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

		if flags.Changed("stakable") {
			stakable, err = flags.GetBool("stakable")
			if err != nil {
				log.Fatal(err)
			}
		}

		if !allFlag {
			if mb != nil && mb.Sharders != nil {
				fmt.Println("MagicBlock Sharders")
				if jsonFlag {
					util.PrettyPrintJSON(mb.Sharders.Nodes)
				} else {
					for _, node := range mb.Sharders.Nodes {
						fmt.Println("ID:", node.ID)
						fmt.Println("  - N2NHost:", node.N2NHost)
						fmt.Println("  - Host:", node.Host)
						fmt.Println("  - Port:", node.Port)
					}
				}
				fmt.Println()
			}
		} else {
			sharders := new(zcncore.MinerSCNodes)

			limit = 20
			offset = 0
			var nodes []zcncore.Node
			for curOff := offset; ; curOff += limit {
				res, err := zcncore.GetSharders(active, stakable, limit, curOff)
				if err != nil {
					log.Fatal(err)
				}

				if err = json.Unmarshal(res, sharders); err != nil {
					log.Fatal(err)
				}

				if len(sharders.Nodes) == 0 {
					break
				}

				nodes = append(nodes, sharders.Nodes...)
			}

			if jsonFlag {
				util.PrettyPrintJSON(nodes)
			} else {
				util.PrintSharderNodes(nodes)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(minerscSharders)
	minerscSharders.PersistentFlags().Bool("json", false, "as JSON")
	minerscSharders.PersistentFlags().Int("limit", 20, "Limits the amount of sharders returned")
	minerscSharders.PersistentFlags().Int("offset", 0, "Skips the number of sharders mentioned")
	minerscSharders.PersistentFlags().Bool("all", false, "include all registered sharders, default returns the first page of sharders")
	minerscSharders.PersistentFlags().Bool("active", true, "Gets active sharders only, set it false to get all sharders")
	minerscSharders.PersistentFlags().Bool("stakable", false, "Gets stakable sharders only if set to true")
}
