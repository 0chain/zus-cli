package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/olekukonko/tablewriter"
)

func PrintJSON(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Failed to convert data to json format : %v", err)
	}
	jsonString := string(b)
	fmt.Println(jsonString)
}

func PrintError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func PrintBlobbers(nodes []*sdk.Blobber, isActive bool) {
	if len(nodes) == 0 {
		if isActive {
			fmt.Println("no active blobbers")
		} else {
			fmt.Println("no blobbers registered yet")
		}
		return
	}
	for _, val := range nodes {
		fmt.Println("- id:                   ", val.ID)
		fmt.Println("  url:                  ", val.BaseURL)
		fmt.Println("  allocated / total capacity:", val.Allocated.String(), "/",
			val.Capacity.String())
		fmt.Println("  last_health_check:	 ", val.LastHealthCheck.ToTime())
		fmt.Println("  terms:")
		fmt.Println("    read_price:         ", val.Terms.ReadPrice.String(), "/ GB")
		fmt.Println("    write_price:        ", val.Terms.WritePrice.String(), "/ GB / time_unit")
		fmt.Println("    max_offer_duration: ", val.Terms.MaxOfferDuration.String())
	}
}

func PrintStakePoolInfo(info *sdk.StakePoolInfo) {
	fmt.Println("pool id:           ", info.ID)
	fmt.Println("balance:           ", info.Balance)
	fmt.Println("total stake:       ", info.StakeTotal)
	fmt.Println("unclaimed rewards: ", info.Rewards)
	fmt.Println("total rewards:     ", info.TotalRewards)
	if len(info.Delegate) == 0 {
		fmt.Println("delegate_pools: no delegate pools")
	} else {
		fmt.Println("delegate_pools:")
		for _, dp := range info.Delegate {
			fmt.Println("- id:               ", dp.ID)
			fmt.Println("  balance:          ", dp.Balance)
			fmt.Println("  delegate_id:      ", dp.DelegateID)
			fmt.Println("  unclaimed reward: ", dp.Rewards)
			fmt.Println("  total_reward:     ", dp.TotalReward)
			fmt.Println("  total_penalty:    ", dp.TotalPenalty)
			fmt.Println("  status:           ", dp.Status)
			fmt.Println("  round_created:    ", dp.RoundCreated)
			fmt.Println("  unstake:          ", dp.UnStake)
			fmt.Println("  staked_at:        ", dp.StakedAt.ToTime().String())
		}
	}
	// settings
	fmt.Println("settings:")
	fmt.Println("  delegate_wallet:  ", info.Settings.DelegateWallet)
	//fmt.Println("  min_stake:        ", info.Settings.MinStake.String())
	//fmt.Println("  max_stake:        ", info.Settings.MaxStake.String())
	fmt.Println("  num_delegates:    ", info.Settings.NumDelegates)
}

func PrintStakePoolUserInfo(info *sdk.StakePoolUserInfo) {
	if len(info.Pools) == 0 {
		fmt.Print("no delegate pools")
		return
	}
	for blobberID, dps := range info.Pools {
		fmt.Println("- blobber_id: ", blobberID)
		for _, dp := range dps {
			fmt.Println("  - id:               ", dp.ID)
			fmt.Println("    balance:          ", dp.Balance)
			fmt.Println("    delegate_id:      ", dp.DelegateID)
			fmt.Println("    unclaimed reward:       ", dp.Rewards)
			fmt.Println("    total rewards:          ", dp.TotalReward)
			fmt.Println("    total penalty:          ", dp.TotalPenalty)
			fmt.Println("    status:          ", dp.Status)
			fmt.Println("    round_created:   ", dp.RoundCreated)
			fmt.Println("    unstake:         ", dp.UnStake)
			fmt.Println("    staked_at:       ", dp.StakedAt.ToTime().String())
		}
	}
}

// WriteTable - Writes string data as a table
func WriteTable(writer io.Writer, header []string, footer []string, data [][]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(header)
	table.SetFooter(footer)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
