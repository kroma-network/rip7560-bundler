package config

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	mapset "github.com/deckarep/golang-set/v2"
)

var (
	EthereumChainID        = big.NewInt(1)
	GoerliChainID          = big.NewInt(5)
	SepoliaChainID         = big.NewInt(11155111)
	OptimismChainID        = big.NewInt(10)
	OptimismGoerliChainID  = big.NewInt(420)
	OptimismSepoliaChainID = big.NewInt(11155420)
	BaseChainID            = big.NewInt(8453)
	BaseGoerliChainID      = big.NewInt(84531)
	BaseSepoliaChainID     = big.NewInt(84532)

	OpStackChains = mapset.NewSet(
		OptimismChainID.Uint64(),
		OptimismGoerliChainID.Uint64(),
		OptimismSepoliaChainID.Uint64(),
		BaseChainID.Uint64(),
		BaseGoerliChainID.Uint64(),
		BaseSepoliaChainID.Uint64(),
	)

	EntryPointAddress     = common.HexToAddress("0x0000000000000000000000000000000000007560")
	NonceManagerAddress   = common.HexToAddress("0x4200000000000000000000000000000000000024")
	DeployerCallerAddress = common.HexToAddress("0x00000000000000000000000000000000ffff7560")
)
