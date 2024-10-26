package util

import (
	"fmt"
	"os"
	"strings"
)

// SplitArgs splt args into slice, and remove if item is empty
func SplitArgs(args string) []string {
	items := strings.Split(args, " ")

	list := make([]string, 0, len(items))

	for _, it := range items {
		trim := strings.TrimSpace(it)
		if len(trim) > 0 {
			list = append(list, trim)
		}
	}

	return list
}

// GetHomeDir Find home directory.
func GetHomeDir() string {
	// Find home directory.
	idr, err := os.UserHomeDir()
	if err != nil {
		PrintError(err)
		os.Exit(1)
	}

	return idr
}

func ExitWithError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}
