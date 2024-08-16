package testutils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"math/big"
)

var (
	OneETH             = big.NewInt(1000000000000000000)
	ValidAddress1      = common.HexToAddress("0x7357b8a705328FC283dF72D7Ac546895B596DC12")
	ValidAddress2      = common.HexToAddress("0x7357c9504B8686c008CCcD6ea47f1c21B7475dE3")
	ValidAddress3      = common.HexToAddress("0x7357C8D931e8cde8ea1b777Cf8578f4A7071f100")
	ValidAddress4      = common.HexToAddress("0x73574a159D05d20FF50D5504057D5C86f2d02a45")
	ValidAddress5      = common.HexToAddress("0x7357C1Fc72a14399cb845f2f71421B4CE7eCE608")
	DummyDeployerData  = []byte{0x12, 0x34, 0xab, 0xcd}
	DummyPaymasterData = []byte{0x12, 0x34, 0xab, 0xcd}
	DummyNonce0        = uint64(0)
	DummyNonce1        = uint64(1)
	ChainID            = big.NewInt(1)
	pk, _              = crypto.GenerateKey()
	DummyEOA, _        = signer.New(hexutil.Encode(crypto.FromECDSA(pk))[2:])
	MockHash           = "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddeaddead"
)
