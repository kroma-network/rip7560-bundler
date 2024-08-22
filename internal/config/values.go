package config

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/entities"
)

type Values struct {
	// Documented variables.
	PrivateKey          string
	EthClientUrl        string
	Port                int
	DataDirectory       string
	MaxVerificationGas  *big.Int
	MaxBatchGasLimit    *big.Int
	MaxTxTTL            time.Duration
	ReputationConstants *entities.ReputationConstants

	// Searcher mode variables.
	EthBuilderUrls    []string
	BlocksInTheFuture int

	// Undocumented variables.
	DebugMode bool
	GinMode   string
}

func envKeyValStringToMap(s string) map[string]string {
	out := map[string]string{}
	for _, pair := range strings.Split(s, "&") {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			break
		}
		out[kv[0]] = kv[1]
	}
	return out
}

func envArrayToAddressSlice(s string) []common.Address {
	env := strings.Split(s, ",")
	slc := []common.Address{}
	for _, ep := range env {
		slc = append(slc, common.HexToAddress(strings.TrimSpace(ep)))
	}

	return slc
}

func envArrayToStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func variableNotSetOrIsNil(env string) bool {
	return !viper.IsSet(env) || viper.GetString(env) == ""
}

// GetValues returns config for the bundler that has been read in from env vars. See
// https://docs.stackup.sh/docs/packages/bundler/configure for details.
func GetValues() *Values {
	// Default variables
	viper.SetDefault("rip7560_bundler_port", 7560)
	viper.SetDefault("rip7560_bundler_data_directory", "/tmp/stackup_bundler")
	viper.SetDefault("rip7560_bundler_supported_entry_points", "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
	viper.SetDefault("rip7560_bundler_max_verification_gas", 6000000)
	// TODO : adjust args from geth request, deprecate this!
	viper.SetDefault("rip7560_bundler_max_batch_gas_limit", 18000000)
	viper.SetDefault("rip7560_bundler_max_tx_ttl_seconds", 180)
	viper.SetDefault("rip7560_bundler_blocks_in_the_future", 6) // TODO!
	viper.SetDefault("rip7560_bundler_debug_mode", false)
	viper.SetDefault("rip7560_bundler_gin_mode", gin.ReleaseMode)

	// Read in from .env file if available
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			// Can ignore
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	// Read in from environment variables
	_ = viper.BindEnv("rip7560_bundler_eth_client_url")
	_ = viper.BindEnv("rip7560_bundler_private_key")
	_ = viper.BindEnv("rip7560_bundler_port")
	_ = viper.BindEnv("rip7560_bundler_data_directory")
	_ = viper.BindEnv("rip7560_bundler_supported_entry_points")
	_ = viper.BindEnv("rip7560_bundler_max_verification_gas")
	_ = viper.BindEnv("rip7560_bundler_max_batch_gas_limit")
	_ = viper.BindEnv("rip7560_bundler_max_tx_ttl_seconds")
	_ = viper.BindEnv("rip7560_bundler_eth_builder_urls")
	_ = viper.BindEnv("rip7560_bundler_blocks_in_the_future")
	_ = viper.BindEnv("rip7560_bundler_debug_mode")
	_ = viper.BindEnv("rip7560_bundler_gin_mode")

	// Validate required variables
	if variableNotSetOrIsNil("rip7560_bundler_eth_client_url") {
		panic("Fatal config error: rip7560_bundler_eth_client_url not set")
	}

	if variableNotSetOrIsNil("rip7560_bundler_private_key") {
		panic("Fatal config error: rip7560_bundler_private_key not set")
	}

	switch viper.GetString("mode") {
	case "searcher":
		if variableNotSetOrIsNil("rip7560_bundler_eth_builder_urls") {
			panic("Fatal config error: rip7560_bundler_eth_builder_urls not set")
		}
	}

	// Return Values
	privateKey := viper.GetString("rip7560_bundler_private_key")
	ethClientUrl := viper.GetString("rip7560_bundler_eth_client_url")
	port := viper.GetInt("rip7560_bundler_port")
	dataDirectory := viper.GetString("rip7560_bundler_data_directory")
	maxVerificationGas := big.NewInt(int64(viper.GetInt("rip7560_bundler_max_verification_gas")))
	maxBatchGasLimit := big.NewInt(int64(viper.GetInt("rip7560_bundler_max_batch_gas_limit")))
	maxTxTTL := time.Second * viper.GetDuration("rip7560_bundler_max_tx_ttl_seconds")
	ethBuilderUrls := envArrayToStringSlice(viper.GetString("rip7560_bundler_eth_builder_urls"))
	blocksInTheFuture := viper.GetInt("rip7560_bundler_blocks_in_the_future")
	debugMode := viper.GetBool("rip7560_bundler_debug_mode")
	ginMode := viper.GetString("rip7560_bundler_gin_mode")
	return &Values{
		PrivateKey:          privateKey,
		EthClientUrl:        ethClientUrl,
		Port:                port,
		DataDirectory:       dataDirectory,
		MaxVerificationGas:  maxVerificationGas,
		MaxBatchGasLimit:    maxBatchGasLimit,
		MaxTxTTL:            maxTxTTL,
		ReputationConstants: NewReputationConstantsFromEnv(),
		EthBuilderUrls:      ethBuilderUrls,
		BlocksInTheFuture:   blocksInTheFuture,
		DebugMode:           debugMode,
		GinMode:             ginMode,
	}
}
