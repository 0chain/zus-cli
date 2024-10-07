package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

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

// WriteTable - Writes string data as a table
func WriteTable(writer io.Writer, header []string, footer []string, data [][]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(header)
	table.SetFooter(footer)
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()
}
