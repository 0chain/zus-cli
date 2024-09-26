package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "zus-cli",
	// TODO: update Short and Long description
	Short: "zus-cli is a decentralized storage application written on the 0Chain platform",
	Long: `zus-cli is a decentralized storage application written on the 0Chain platform.
			Complete documentation is available at https://docs.zus.network/guides/zbox-cli`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
