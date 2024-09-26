package main

import "github.com/0chain/zus-cli/cmd"

var VersionStr string

func main() {
	cmd.VersionStr = VersionStr
	cmd.Execute()
}
