// Package client provides the mediator for processing incoming UserOperations to the bundler.
package client

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-logr/logr"
	"github.com/stackup-wallet/stackup-bundler/internal/logger"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/filter"
	"github.com/stackup-wallet/stackup-bundler/pkg/entrypoint/stake"
	"github.com/stackup-wallet/stackup-bundler/pkg/gas"
	"github.com/stackup-wallet/stackup-bundler/pkg/mempool"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules"
	"github.com/stackup-wallet/stackup-bundler/pkg/modules/noop"
	"github.com/stackup-wallet/stackup-bundler/pkg/state"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// Client controls the end to end process of adding incoming UserOperations to the mempool. It also
// implements the required RPC methods as specified in EIP-4337.
type Client struct {
	mempool                 *mempool.Mempool
	chainID                 *big.Int
	supportedEntryPoints    []common.Address
	userOpHandler           modules.UserOpHandlerFunc
	logger                  logr.Logger
	getUserOpReceipt        GetUserOpReceiptFunc
	getRip7560UserOpReceipt GetRip7560UserOpReceiptFunc
	getGasPrices            GetGasPricesFunc
	getGasEstimate          GetGasEstimateFunc
	getUserOpByHash         GetUserOpByHashFunc
	getStakeFunc            stake.GetStakeFunc
	opLookupLimit           uint64
}

// New initializes a new ERC-4337 client which can be extended with modules for validating UserOperations
// that are allowed to be added to the mempool.
func New(
	mempool *mempool.Mempool,
	chainID *big.Int,
	supportedEntryPoints []common.Address,
	opLookupLimit uint64,
) *Client {
	return &Client{
		mempool:                 mempool,
		chainID:                 chainID,
		supportedEntryPoints:    supportedEntryPoints,
		userOpHandler:           noop.UserOpHandler,
		logger:                  logger.NewZeroLogr().WithName("client"),
		getUserOpReceipt:        getUserOpReceiptNoop(),
		getRip7560UserOpReceipt: getRip7560UserOpReceiptNoop(),
		getGasPrices:            getGasPricesNoop(),
		getGasEstimate:          getGasEstimateNoop(),
		getUserOpByHash:         getUserOpByHashNoop(),
		getStakeFunc:            stake.GetStakeFuncNoop(),
		opLookupLimit:           opLookupLimit,
	}
}

func (i *Client) parseEntryPointAddress(ep string) (common.Address, error) {
	for _, addr := range i.supportedEntryPoints {
		if common.HexToAddress(ep) == addr {
			return addr, nil
		}
	}

	return common.Address{}, errors.New("entryPoint: Implementation not supported")
}

// UseLogger defines the logger object used by the Client instance based on the go-logr/logr interface.
func (i *Client) UseLogger(logger logr.Logger) {
	i.logger = logger.WithName("client")
}

// UseModules defines the UserOpHandlers to process a userOp after it has gone through the standard checks.
func (i *Client) UseModules(handlers ...modules.UserOpHandlerFunc) {
	i.userOpHandler = modules.ComposeUserOpHandlerFunc(handlers...)
}

// SetGetUserOpReceiptFunc defines a general function for fetching a UserOpReceipt given a userOpHash and
// EntryPoint address. This function is called in *Client.GetUserOperationReceipt.
func (i *Client) SetGetUserOpReceiptFunc(fn GetUserOpReceiptFunc) {
	i.getUserOpReceipt = fn
}

// SetGetUserOpRip7560ReceiptFunc defines a general function for fetching a UserOpReceipt given a userOpHash and
// EntryPoint address. This function is called in *Client.GetUserOperationReceipt.
func (i *Client) SetGetRip7560UserOpReceiptFunc(fn GetRip7560UserOpReceiptFunc) {
	i.getRip7560UserOpReceipt = fn
}

// SetGetGasPricesFunc defines a general function for fetching values for maxFeePerGas and
// maxPriorityFeePerGas. This function is called in *Client.EstimateUserOperationGas if given fee values are
// 0.
func (i *Client) SetGetGasPricesFunc(fn GetGasPricesFunc) {
	i.getGasPrices = fn
}

// SetGetGasEstimateFunc defines a general function for fetching an estimate for verificationGasLimit and
// callGasLimit given a userOp and EntryPoint address. This function is called in
// *Client.EstimateUserOperationGas.
func (i *Client) SetGetGasEstimateFunc(fn GetGasEstimateFunc) {
	i.getGasEstimate = fn
}

// SetGetUserOpByHashFunc defines a general function for fetching a userOp given a userOpHash, EntryPoint
// address, and chain ID. This function is called in *Client.GetUserOperationByHash.
func (i *Client) SetGetUserOpByHashFunc(fn GetUserOpByHashFunc) {
	i.getUserOpByHash = fn
}

// SetGetStakeFunc defines a general function for retrieving the EntryPoint stake for a given address. This
// function is called in *Client.SendUserOperation to create a context.
func (i *Client) SetGetStakeFunc(fn stake.GetStakeFunc) {
	i.getStakeFunc = fn
}

// SendUserOperation implements the method call for eth_sendUserOperation.
// It returns true if userOp was accepted otherwise returns an error.
func (i *Client) SendUserOperation(op map[string]any, ep string) (string, error) {
	// Init logger
	l := i.logger.WithName("eth_sendUserOperation")
	epAddr := common.Address{}
	l = l.WithValues("chain_id", i.chainID.String())

	userOp, err := userop.New(op)
	if err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}

	txHash, err := getTransactionHashByUserOp(op)
	if err != nil {
		l.Error(err, "getTransactionHashByUserOp error")
		return "", err
	}
	l = l.WithValues("txHash", txHash)

	// Run through client module stack.
	ctx, err := modules.NewUserOpHandlerContext(
		userOp,
		epAddr,
		i.chainID,
		i.mempool,
		i.getStakeFunc,
	)
	if err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}
	if err := i.userOpHandler(ctx); err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}

	// Add userOp to mempool.
	if err := i.mempool.AddOp(epAddr, ctx.UserOp); err != nil {
		l.Error(err, "eth_sendUserOperation error")
		return "", err
	}

	l.Info("eth_sendUserOperation ok")
	return txHash.Hex(), nil
}

// EstimateUserOperationGas returns estimates for PreVerificationGas, VerificationGasLimit, and CallGasLimit
// given a UserOperation, EntryPoint address, and state OverrideSet. The signature field and current gas
// values will not be validated although there should be dummy values in place for the most reliable results
// (e.g. a signature with the correct length).
func (i *Client) EstimateUserOperationGas(
	op map[string]any,
	ep string,
	os map[string]any,
) (*gas.GasEstimates, error) {
	// Init logger
	l := i.logger.WithName("eth_estimateUserOperationGas")

	// Check EntryPoint and userOp is valid.
	//epAddr, err := i.parseEntryPointAddress(ep)
	//if err != nil {
	//	l.Error(err, "eth_estimateUserOperationGas error")
	//	return nil, err
	//}
	//l = l.
	//	WithValues("entrypoint", epAddr.String()).
	//	WithValues("chain_id", i.chainID.String())

	userOp, err := userop.New(op)
	if err != nil {
		l.Error(err, "eth_estimateUserOperationGas error")
		return nil, err
	}
	//hash := userOp.GetUserOpHash(epAddr, i.chainID)
	//l = l.WithValues("userop_hash", hash)

	// Parse state override set.
	sos, err := state.ParseOverrideData(os)
	if err != nil {
		l.Error(err, "eth_estimateUserOperationGas error")
		return nil, err
	}

	// Override op with suggested gas prices if maxFeePerGas is 0. This allows for more reliable gas
	// estimations upstream. The default balance override also ensures simulations won't revert on
	// insufficient funds.
	if userOp.MaxFeePerGas.Cmp(common.Big0) != 1 {
		gp, err := i.getGasPrices()
		if err != nil {
			l.Error(err, "eth_estimateUserOperationGas error")
			return nil, err
		}
		userOp.MaxFeePerGas = gp.MaxFeePerGas
		userOp.MaxPriorityFeePerGas = gp.MaxPriorityFeePerGas
	}

	// Estimate gas limits
	vg, cg, err := i.getGasEstimate(userOp, sos)
	if err != nil {
		l.Error(err, "eth_estimateUserOperationGas error")
		return nil, err
	}

	// Calculate PreVerificationGas
	//pvg, err := i.ov.CalcPreVerificationGasWithBuffer(userOp)
	//if err != nil {
	//	l.Error(err, "eth_estimateUserOperationGas error")
	//	return nil, err
	//}

	l.Info("eth_estimateUserOperationGas ok")
	return &gas.GasEstimates{
		PreVerificationGas:   big.NewInt(0), // do not use this
		VerificationGasLimit: big.NewInt(int64(vg)),
		CallGasLimit:         big.NewInt(int64(cg)),

		// TODO: Deprecate in v0.7
		VerificationGas: big.NewInt(int64(vg)),
	}, nil
}

// GetUserOperationReceipt fetches a UserOperation receipt based on a userOpHash returned by
// *Client.SendUserOperation.
func (i *Client) GetUserOperationReceipt(
	hash string,
) (*filter.UserOperationReceipt, error) {
	// Init logger
	l := i.logger.WithName("eth_getUserOperationReceipt").WithValues("userop_hash", hash)

	ev, err := i.getUserOpReceipt(hash, i.supportedEntryPoints[0], i.opLookupLimit)
	if err != nil {
		l.Error(err, "eth_getUserOperationReceipt error")
		return nil, err
	}

	l.Info("eth_getUserOperationReceipt ok")
	return ev, nil
}

// GetRIP7560UserOperationReceipt returns RIP7560 transaction receipt based on a userOp
// *Client.SendUserOperation.
func (i *Client) GetRIP7560UserOperationReceipt(
	op userOperation,
) (*types.Receipt, error) {
	// Init logger
	l := i.logger.WithName("eth_getUserOperationRip7560Receipt").WithValues("rip7560userop")

	txHash, err := getTransactionHashByUserOp(op)
	if err != nil {
		l.Error(err, "getTransactionHashByUserOp error")
		return nil, err
	}

	receipt, err := i.getRip7560UserOpReceipt(txHash, i.opLookupLimit)
	if err != nil {
		l.Error(err, "getRip7560UserOpReceipt error")
	}

	l.Info("eth_getUserOperationReceipt ok")
	return receipt, nil
}

// GetUserOperationByHash returns a UserOperation based on a given userOpHash returned by
// *Client.SendUserOperation.
func (i *Client) GetUserOperationByHash(hash string) (*filter.HashLookupResult, error) {
	// Init logger
	l := i.logger.WithName("eth_getUserOperationByHash").WithValues("userop_hash", hash)

	res, err := i.getUserOpByHash(hash, i.supportedEntryPoints[0], i.chainID, i.opLookupLimit)
	if err != nil {
		l.Error(err, "eth_getUserOperationByHash error")
		return nil, err
	}

	return res, nil
}

// SupportedEntryPoints implements the method call for eth_supportedEntryPoints. It returns the array of
// EntryPoint addresses that is supported by the client. The first address in the array is the preferred
// EntryPoint.
func (i *Client) SupportedEntryPoints() ([]string, error) {
	slc := []string{}
	for _, ep := range i.supportedEntryPoints {
		slc = append(slc, ep.String())
	}

	return slc, nil
}

// ChainID implements the method call for eth_chainId. It returns the current chainID used by the client.
// This method is used to validate that the client's chainID is in sync with the caller.
func (i *Client) ChainID() (string, error) {
	return hexutil.EncodeBig(i.chainID), nil
}

func getTransactionHashByUserOp(op userOperation) (common.Hash, error) {
	userOp, err := userop.New(op)
	if err != nil {
		//l.Error(err, "eth_sendUserOperation error")
		return common.Hash{}, err
	}

	txArgs := CreateUserOperationArgs(userOp)
	gas, _ := strconv.ParseUint(txArgs.Gas[2:], 16, 64)
	sender := common.HexToAddress(txArgs.Sender)
	builderFee, _ := new(big.Int).SetString(txArgs.BuilderFee, 0)
	validationGas, _ := strconv.ParseUint(txArgs.ValidationGas[2:], 16, 64)
	paymasterGas, _ := strconv.ParseUint(txArgs.PaymasterGas[2:], 16, 64)
	postopGas, _ := strconv.ParseUint(txArgs.PostOpGas[2:], 16, 64)
	bigNonce, _ := new(big.Int).SetString(txArgs.BigNonce, 0)
	aatx := types.Rip7560AccountAbstractionTx{
		To:         &common.Address{},
		ChainID:    nil,
		Gas:        gas,
		GasTipCap:  userOp.MaxPriorityFeePerGas,
		GasFeeCap:  userOp.MaxFeePerGas,
		Value:      nil,
		Data:       userOp.CallData,
		AccessList: types.AccessList{},
		// RIP-7560 parameters
		Sender:        &sender,
		Signature:     userOp.Signature,
		PaymasterData: userOp.PaymasterAndData,
		DeployerData:  userOp.InitCode,
		BuilderFee:    builderFee,
		ValidationGas: validationGas,
		PaymasterGas:  paymasterGas,
		PostOpGas:     postopGas,
		// RIP-7712 parameter
		BigNonce: bigNonce,
	}
	tx := types.NewTx(&aatx)
	return tx.Hash(), nil
}
