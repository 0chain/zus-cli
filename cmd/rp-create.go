package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/spf13/cobra"
)

// rpCreate creates read pool
var rpCreate = &cobra.Command{
	Use:   "rp-create",
	Short: "Create read pool if missing",
	Long:  `Create read pool in storage SC if the pool is missing.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if _, _, err = sdk.CreateReadPool(); err != nil {
			log.Fatalf("Failed to create read pool: %v\n", err)
		}
		fmt.Println("Read pool created successfully")
	},
}

func init() {
	rootCmd.AddCommand(rpCreate)
}
