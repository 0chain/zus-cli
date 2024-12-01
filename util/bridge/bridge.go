package bridge

import (
	"fmt"
	"path/filepath"

	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zus-cli/util"
	"github.com/spf13/cobra"
)

const (
	DefaultRetries = 60
)

const (
	DefaultConfigChainFileName = "config.yaml"
	DefaultWalletFileName      = "wallet.json"
)

const (
	OptionHash            = "hash"          // OptionHash hash passed to cmd
	OptionAmount          = "amount"        // OptionAmount amount passed to cmd
	OptionToken           = "token"         // OptionToken token in SAS passed to cmd
	OptionRetries         = "retries"       // OptionRetries retries
	OptionConfigFolder    = "path"          // OptionConfigFolder config folder
	OptionChainConfigFile = "chain_config"  // OptionChainConfigFile sdk config filename
	OptionMnemonic        = "mnemonic"      // OptionMnemonic bridge config filename
	OptionKeyPassword     = "password"      // OptionKeyPassword bridge config filename
	OptionAccountIndex    = "account_index" // OptionAccountIndex ethereum account index
	OptionAddressIndex    = "address_index" // OptionAddressIndex ethereum address index
	OptionBip32           = "bip32"         // OptionBip32 use bip32 derivation path
	OptionClientKey       = "client_key"
	OptionClientID        = "client_id"
	OptionEthereumAddress = "ethereum_address"
	OptionURL             = "url"
	OptionMinStake        = "min_stake"
	OptionMaxStake        = "max_stake"
	OptionNumDelegates    = "num_delegates"
	OptionServiceCharge   = "service_charge"
	OptionWalletFile      = "wallet"
)

type CommandWithBridge func(*zcnbridge.BridgeClient, ...*Arg)
type Command func(...*Arg)

type Option struct {
	Name         string
	Value        interface{}
	Typename     string
	Usage        string
	MissingError string
	Required     bool
}

type Arg struct {
	typeName  string
	fieldName string
	value     interface{}
}

var (
	walletFileOption = &Option{
		Name:         OptionWalletFile,
		Value:        "wallet.json",
		Typename:     "string",
		Usage:        "Wallet file",
		MissingError: "Wallet file not specified",
		Required:     false,
	}

	configFolderOption = &Option{
		Name:         OptionConfigFolder,
		Value:        util.GetDefaultConfigDirOrExit(),
		Typename:     "string",
		Usage:        "Config home folder",
		MissingError: "Config home folder not specified",
		Required:     false,
	}

	configChainFileOption = &Option{
		Name:         OptionChainConfigFile,
		Value:        DefaultConfigChainFileName,
		Typename:     "string",
		Usage:        "Chain config file name",
		MissingError: "Chain config file name not specified",
		Required:     false,
	}
)

func WithRetries(usage string) *Option {
	return &Option{
		Name:         OptionRetries,
		Value:        DefaultRetries,
		Typename:     "int",
		Usage:        usage,
		MissingError: "Retries count should be provided",
		Required:     false,
	}
}

func WithToken(usage string) *Option {
	return &Option{
		Name:         OptionToken,
		Value:        float64(0),
		Usage:        usage,
		Typename:     "float64",
		MissingError: "Token should be provided",
		Required:     true,
	}
}

func WithAmount(usage string) *Option {
	return &Option{
		Name:         OptionAmount,
		Value:        int64(0),
		Usage:        usage,
		Typename:     "int64",
		MissingError: "Amount should be provided",
		Required:     true,
	}
}

func WithHash(usage string) *Option {
	return &Option{
		Name:         OptionHash,
		Value:        "",
		Usage:        usage,
		Typename:     "string",
		MissingError: "hash of the transaction should be provided",
		Required:     true,
	}
}

func GetChainConfigFile(args []*Arg) string {
	return GetString(args, OptionChainConfigFile)
}

func GetConfigFolder(args []*Arg) string {
	return GetString(args, OptionConfigFolder)
}

func GetHash(args []*Arg) string {
	return GetString(args, OptionHash)
}

func GetAmount(args []*Arg) uint64 {
	return uint64(GetInt64(args, OptionAmount))
}

func GetToken(args []*Arg) float64 {
	return GetFloat64(args, OptionToken)
}

func GetRetries(args []*Arg) int {
	return GetInt(args, OptionRetries)
}

func GetClientID(args []*Arg) string {
	return GetString(args, OptionClientID)
}

func GetClientKey(args []*Arg) string {
	return GetString(args, OptionClientKey)
}

func GetEthereumAddress(args []*Arg) string {
	return GetString(args, OptionEthereumAddress)
}

func GetURL(args []*Arg) string {
	return GetString(args, OptionURL)
}

func GetMinStake(args []*Arg) int64 {
	return GetInt64(args, OptionMinStake)
}

func GetMaxStake(args []*Arg) int64 {
	return GetInt64(args, OptionMaxStake)
}

func GetNumDelegates(args []*Arg) int {
	return GetInt(args, OptionNumDelegates)
}

func GetServiceCharge(args []*Arg) float64 {
	return GetFloat64(args, OptionServiceCharge)
}

func GetWalletFile(args []*Arg) string {
	return GetString(args, OptionWalletFile)
}

func GetString(args []*Arg, name string) string {
	if len(args) == 0 {
		util.ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(string)
		}
	}

	util.ExitWithError("failed to get " + name)

	return ""
}

func GetInt(args []*Arg, name string) int {
	if len(args) == 0 {
		util.ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(int)
		}
	}

	util.ExitWithError("failed to get " + name)

	return 0
}

func GetFloat64(args []*Arg, name string) float64 {
	if len(args) == 0 {
		util.ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(float64)
		}
	}

	util.ExitWithError("failed to get " + name)

	return 0
}

func GetInt64(args []*Arg, name string) int64 {
	if len(args) == 0 {
		util.ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(int64)
		}
	}

	util.ExitWithError("failed to get " + name)

	return 0
}

func GetUint64(args []*Arg, name string) uint64 {
	if len(args) == 0 {
		util.ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(uint64)
		}
	}

	util.ExitWithError("failed to get " + name)

	return 0
}

// CreateCommand Function to initialize bridge commands with DRY principle
func CreateCommand(use, short, long string, functor Command, hidden bool, opts ...*Option) *cobra.Command {
	fn := func(parameters ...*Arg) {
		functor(parameters...)
	}

	opts = append(opts, configFolderOption)
	opts = append(opts, configChainFileOption)

	command := createBridgeComm(use, short, long, fn, opts, hidden)
	AppendOptions(opts, command)
	return command
}

// CreateCommandWithBridge Function to initialize bridge commands with DRY principle
func CreateCommandWithBridge(use, short, long string, functor CommandWithBridge, hidden bool, opts ...*Option) *cobra.Command {
	fn := func(parameters ...*Arg) {
		folder := GetConfigFolder(parameters)
		chainConfigFile := GetChainConfigFile(parameters)

		bridge := initBridge(folder, chainConfigFile)
		functor(bridge, parameters...)
	}

	opts = append(opts, configFolderOption)
	opts = append(opts, configChainFileOption)
	command := createBridgeComm(use, short, long, fn, opts, hidden)
	AppendOptions(opts, command)
	return command
}

func AppendOptions(opts []*Option, command *cobra.Command) {
	for _, opt := range opts {
		switch opt.Typename {
		case "string":
			command.PersistentFlags().String(opt.Name, opt.Value.(string), opt.Usage)
		case "int64":
			command.PersistentFlags().Int64(opt.Name, opt.Value.(int64), opt.Usage)
		case "float64":
			command.PersistentFlags().Float64(opt.Name, opt.Value.(float64), opt.Usage)
		case "int":
			command.PersistentFlags().Int(opt.Name, opt.Value.(int), opt.Usage)
		}

		if opt.Required {
			_ = command.MarkFlagRequired(opt.Name)
		}
	}
}

func createBridgeComm(
	use string,
	short string,
	long string,
	functor Command,
	opts []*Option,
	hidden bool,
) *cobra.Command {
	var cobraCommand = &cobra.Command{
		Use:    use,
		Short:  short,
		Long:   long,
		Args:   cobra.MinimumNArgs(0),
		Hidden: hidden,
		Run: func(cmd *cobra.Command, args []string) {
			fflags := cmd.Flags()

			var parameters []*Arg

			for _, opt := range opts {
				if !fflags.Changed(opt.Name) && opt.Required {
					//TODO: add default missing error
					util.ExitWithError(opt.MissingError)
				}

				var arg *Arg
				switch opt.Typename {
				case "string":
					optValue, err := fflags.GetString(opt.Name)
					if err != nil {
						util.ExitWithError(err)
					}
					arg = &Arg{
						typeName:  opt.Typename,
						fieldName: opt.Name,
						value:     optValue,
					}
				case "int64":
					optValue, err := fflags.GetInt64(opt.Name)
					if err != nil {
						util.ExitWithError(err)
					}
					arg = &Arg{
						typeName:  opt.Typename,
						fieldName: opt.Name,
						value:     optValue,
					}
				case "float64":
					optValue, err := fflags.GetFloat64(opt.Name)
					if err != nil {
						util.ExitWithError(err)
					}
					arg = &Arg{
						typeName:  opt.Typename,
						fieldName: opt.Name,
						value:     optValue,
					}
				case "int":
					optValue, err := fflags.GetInt(opt.Name)
					if err != nil {
						util.ExitWithError(err)
					}
					arg = &Arg{
						typeName:  opt.Typename,
						fieldName: opt.Name,
						value:     optValue,
					}
				default:
					util.ExitWithError(fmt.Printf("unknown argument: %s, value: %v\n", opt.Name, opt.Value))
				}

				parameters = append(parameters, arg)
			}

			// check SDK EthereumNode
			clientConfig, _ := conf.GetClientConfig()
			if clientConfig.EthereumNode == "" {
				util.ExitWithError("ethereum_node_url must be setup in config")
			}

			functor(parameters...)
		},
	}
	return cobraCommand
}

func initBridge(overrideConfigFolder, overrideConfigFile string) *zcnbridge.BridgeClient {
	var (
		configDir           = util.GetDefaultConfigDirOrExit()
		configChainFileName = DefaultConfigChainFileName
		logPath             = "logs"
		loglevel            = "info"
		development         = false
	)

	if overrideConfigFolder != "" {
		configDir = overrideConfigFolder
	}

	configChainFileName = overrideConfigFile

	configDir, err := filepath.Abs(configDir)
	if err != nil {
		util.ExitWithError(err)
	}

	cfg := &zcnbridge.BridgeSDKConfig{
		ConfigDir:       &configDir,
		ConfigChainFile: &configChainFileName,
		LogPath:         &logPath,
		LogLevel:        &loglevel,
		Development:     &development,
	}

	bridge := zcnbridge.SetupBridgeClientSDK(cfg)

	return bridge
}

func Check(cmd *cobra.Command, flags ...string) {
	fflags := cmd.Flags()
	for _, flag := range flags {
		if !fflags.Changed(flag) {
			util.ExitWithError(fmt.Sprintf("Error: '%s' flag is missing", flag))
		}
	}
}
