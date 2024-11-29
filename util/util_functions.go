package util

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/0chain/gosdk/core/conf"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// default configuration
//
//go:embed config.yaml
var ConfigStr string

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

// returns full path of application's default configuration directory
func GetDefaultConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appConfigDir := userConfigDir + string(os.PathSeparator) + ".zcn"
	return appConfigDir, nil
}

// loads and returns default configuration
func LoadDefaultConfig() (conf.Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(strings.NewReader(ConfigStr))
	if err != nil {
		fmt.Println("error reading default config:", err)
		return conf.Config{}, err
	}
	cfg, err := conf.LoadConfig(v)
	if err != nil {
		fmt.Println("error loading default config:", err)
		return conf.Config{}, err
	}
	return cfg, nil
}

func ExitWithError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func SetupInputMap(flags *pflag.FlagSet, sKeys, sValues string) map[string]string {
	var err error
	var keys []string
	if flags.Changed(sKeys) {
		keys, err = flags.GetStringSlice(sKeys)
		if err != nil {
			log.Fatal(err)
		}
	}

	var values []string
	if flags.Changed(sValues) {
		values, err = flags.GetStringSlice(sValues)
		if err != nil {
			log.Fatal(err)
		}
	}

	input := make(map[string]string)
	if len(keys) != len(values) {
		log.Fatal("number " + sKeys + " must equal the number " + sValues)
	}
	for i := 0; i < len(keys); i++ {
		v := strings.TrimSpace(values[i])
		k := strings.TrimSpace(keys[i])
		input[k] = v
	}
	return input
}
