package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

// rpInfo information
var rpInfo = &cobra.Command{
	Use:   "rp-info",
	Short: "Read pool information.",
	Long:  `Read pool information.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		doJSON, _ := cmd.Flags().GetBool("json")

		info, err := sdk.GetReadPoolInfo("")
		if err != nil {
			log.Fatalf("Failed to get read pool info: %v", err)
		}

		token, err := info.Balance.ToToken()
		if err != nil {
			log.Fatal(err)
		}
		usd, err := zcncore.ConvertTokenToUSD(token)
		var bt = float64(info.Balance) / 1e10
		if err != nil {
			log.Fatalf("Failed to convert token to usd: %v", err)
		}

		if info.Balance == 0 {
			fmt.Println("no tokens locked")
			return
		}

		if doJSON {
			jsonCurrencies := map[string]interface{}{
				"usd": usd,
				"zcn": bt,
				"fmt": info.Balance,
			}

			util.PrintJSON(jsonCurrencies)
			return
		}
		fmt.Printf("\nRead pool Balance: %v (%.2f USD)\n", info.Balance, usd)
	},
}

func init() {
	rootCmd.AddCommand(rpInfo)
	rpInfo.Flags().Bool("json", false, "(default false) pass this option to print response as json data")
}
