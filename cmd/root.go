package cmd

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/0chain/gosdk/core/client"
	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/core/logger"
	"github.com/0chain/gosdk/core/zcncrypto"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd flags
var cfgFile string
var networkFile string
var walletFile string
var walletClientID string
var walletClientKey string
var cDir string
var nonce int64
var txFee float64
var bSilent bool

// default configuration
//
//go:embed config.yaml
var configStr string

// wallet info
var walletJSON string
var clientWallet *zcncrypto.Wallet

var cfg conf.Config
var network conf.Network

var rootCmd = &cobra.Command{
	Use: "zus-cli",
	// TODO: update Short and Long description
	Short: "zus-cli is a decentralized storage application written on the 0Chain platform",
	Long: `zus-cli is a decentralized storage application written on the 0Chain platform.
			Complete documentation is available at https://docs.zus.network/guides/zbox-cli`,
}

func init() {

	InstallDLLs()

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().StringVar(&networkFile, "network", "network.yaml", "network file to overwrite the network details (if required, default is network.yaml)")
	rootCmd.PersistentFlags().StringVar(&walletFile, "wallet", "wallet.json", "wallet file (default is wallet.json)")
	rootCmd.PersistentFlags().StringVar(&walletClientID, "wallet_client_id", "", "wallet client_id")
	rootCmd.PersistentFlags().StringVar(&walletClientKey, "wallet_client_key", "", "wallet client_key")
	rootCmd.PersistentFlags().Int64Var(&nonce, "withNonce", 0, "nonce that will be used in transaction (default is 0)")
	rootCmd.PersistentFlags().StringVar(&cDir, "configDir", "", "configuration directory")
	rootCmd.PersistentFlags().BoolVar(&bSilent, "silent", false, "(default false) Do not show interactive sdk logs (shown by default)")
	rootCmd.PersistentFlags().Float64Var(&txFee, "fee", 0, "transaction fee for the given transaction (if unset, it will be set to blockchain min fee)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// returns full path of application's default configuration directory
func GetAppConfigDir() (string, error) {
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
	err := v.ReadConfig(strings.NewReader(configStr))
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

// returns full path of application configuration directory
func GetConfigDir() (string, error) {
	if len(cDir) != 0 {
		return cDir, nil
	}
	appConfigDir, err := GetAppConfigDir()
	if err != nil {
		return "", err
	}
	return appConfigDir, nil
}

// loads and returns configuration
func LoadConfig() (conf.Config, error) {
	configDir, _ := GetConfigDir()
	cfg, err := conf.LoadConfigFile(filepath.Join(configDir, cfgFile))
	if err == nil {
		return cfg, nil
	}
	fmt.Println("Can't read config:", err)
	fmt.Println("using default config")
	fmt.Printf("config: %v", configStr)
	cfg, err = LoadDefaultConfig()
	if err == nil {
		return cfg, nil
	}
	return conf.Config{}, err
}

// loads and returns network configuration
func LoadNetworkFile() (conf.Network, error) {
	configDir, _ := GetConfigDir()
	network, _ := conf.LoadNetworkFile(filepath.Join(configDir, networkFile))
	return network, nil
}

func initGoSDK() (err error) {
	// syncing loggers
	logger.SyncLoggers([]*logger.Logger{zcncore.GetLogger(), sdk.GetLogger()})

	// set the log file
	zcncore.SetLogFile("cmdlog.log", !bSilent)
	sdk.SetLogFile("cmdlog.log", !bSilent)

	err = client.Init(context.Background(), cfg)
	if err != nil {
		return
	}

	return nil
}

// loads wallet and return
func loadWallet() (*zcncrypto.Wallet, string, error) {
	wallet := &zcncrypto.Wallet{}
	if (len(walletClientID) > 0) && (len(walletClientKey) > 0) {
		wallet.ClientID = walletClientID
		wallet.ClientKey = walletClientKey
		clientBytes, err := json.Marshal(wallet)
		if err != nil {
			err = fmt.Errorf("Invalid wallet data passed:" + walletClientID + " " + walletClientKey)
			return nil, "", err
		}
		return wallet, string(clientBytes), nil
	}
	var walletFilePath string
	if filepath.IsAbs(walletFile) {
		walletFilePath = walletFile
	} else {
		configDir, err := GetConfigDir()
		if err != nil {
			fmt.Println(err.Error())
			return nil, "", err
		}
		walletFilePath = configDir + string(os.PathSeparator) + walletFile
	}

	var walletStr string
	if _, err := os.Stat(walletFilePath); os.IsNotExist(err) {
		walletStr, err = zcncore.CreateWalletOffline()
		if err != nil {
			fmt.Println(err.Error())
			return nil, "", err
		}
		fmt.Println("ZCN wallet created")
		file, err := os.Create(walletFilePath)
		if err != nil {
			fmt.Println(err.Error())
			return nil, "", err
		}
		defer file.Close()
		fmt.Fprint(file, walletStr)
	} else {
		f, err := os.Open(walletFilePath)
		if err != nil {
			fmt.Println("Error opening the wallet", err)
			return nil, "", err
		}
		clientBytes, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("Error reading the wallet", err)
			return nil, "", err
		}
		walletStr = string(clientBytes)
	}
	err := json.Unmarshal([]byte(walletStr), wallet)
	if err != nil {
		fmt.Println("Invalid wallet at path:" + walletFilePath)
		return nil, "", err
	}
	return wallet, walletStr, nil
}

// creates config file "config.yaml" with default configuration in user's configuration directory.
func createConfigFile() error {
	appConfigDir, err := GetAppConfigDir()
	if err != nil {
		fmt.Println("error getting appConfigDir :", err)
		return err
	}
	if err := os.MkdirAll(appConfigDir, 0744); err != nil {
		fmt.Println("error creating default configuration directory :", err)
		return err
	}
	// fmt.Printf("created default configuration directory : %v", appConfigDir)
	configFile := appConfigDir + string(os.PathSeparator) + "config.yaml"

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		file, err := os.Create(configFile)
		if err != nil {
			fmt.Println("error creating configFile :", err)
			return err
		}
		defer file.Close()
		_, err = fmt.Fprint(file, configStr)
		if err != nil {
			fmt.Println("error writing default configuration :", err)
			return err
		}
		fmt.Printf("created default configuration file : %v\n", configFile)
	}
	return nil
}

func initConfig() {

	// DEBUG
	appConfigDir, _ := GetAppConfigDir()
	fmt.Printf("DEBUG: appConfigDir = %v\n", appConfigDir)

	_ = createConfigFile()

	// load config
	var err error
	cfg, err = LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v", err)
		os.Exit(1)
	}

	// load network file
	network, _ = LoadNetworkFile()

	// init gosdk
	err = initGoSDK()
	if err != nil {
		fmt.Println("error initializing SDK: ", err)
		os.Exit(1)
	}

	// load wallet info
	clientWallet, walletJSON, err = loadWallet()
	if err != nil {
		fmt.Println("error loading wallet: ", err)
		os.Exit(1)
	}

	//init the storage sdk with the known miners, sharders and client wallet info
	if err = client.InitSDK(
		walletJSON,
		cfg.BlockWorker,
		cfg.ChainID,
		cfg.SignatureScheme,
		nonce,
		false, true,
		int(zcncore.ConvertToValue(txFee)),
	); err != nil {
		fmt.Println("Error in sdk init", err)
		os.Exit(1)
	}

	sdk.SetNumBlockDownloads(10)
}
