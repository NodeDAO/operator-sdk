// description:
// @author renshiwei
// Date: 2023/6/13

package register

import (
	"context"
	"encoding/hex"
	"github.com/NodeDAO/operator-sdk/common/logger"
	"github.com/NodeDAO/operator-sdk/contracts"
	"github.com/NodeDAO/operator-sdk/eth1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"math/big"
	"strings"
	"time"
)

// RegisterValidator register validator
type RegisterValidator struct {
	// param
	network     string
	maxGasPrice spec.Gwei

	// init
	elClient        *eth1.EthClient
	keyTransactOpts *bind.TransactOpts
	// contracts
	liqContract *contracts.LiqContract
}

// NewRegisterValidator new RegisterValidator
func NewRegisterValidator(ctx context.Context, network, elAddr, controllerAddressPrivateKey string, maxGasPrice spec.Gwei) (*RegisterValidator, error) {
	var err error

	elClient, err := eth1.NewEthClient(ctx, elAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to new eth client. network:%s", network)
	}

	liqContract, err := contracts.NewLiqContract(network, elClient.Client)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to new liqContract. network:%s", network)
	}

	gasPrice := big.NewInt(int64(maxGasPrice))
	if maxGasPrice == 0 {
		gasPrice = nil
	}
	chainID, err := elClient.Client.ChainID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get chainID.")
	}
	opts, err := eth1.KeyTransactOpts(chainID, controllerAddressPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create transact opts.")
	}
	opts.GasPrice = gasPrice

	return &RegisterValidator{
		network:         network,
		elClient:        elClient,
		liqContract:     liqContract,
		maxGasPrice:     maxGasPrice,
		keyTransactOpts: opts,
	}, nil
}

// RegisterCount Calculate the number of validators that can be registered with the operator ID
func (e *RegisterValidator) RegisterCount(operatorId int64) (uint64, error) {
	operatorPoolBalances, err := e.liqContract.Contract.OperatorPoolBalances(nil, big.NewInt(operatorId))
	if err != nil {
		return 0, errors.Wrapf(err, "Fail to get operatorPoolBalances. operatorId:%d", operatorId)
	}

	operatorNftPoolBalances, err := e.liqContract.Contract.OperatorNftPoolBalances(nil, big.NewInt(operatorId))
	if err != nil {
		return 0, errors.Wrapf(err, "Fail to get operatorNftPoolBalances. operatorId:%d", operatorId)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	liqBalance, err := e.elClient.BalanceAt(ctxTimeout, e.liqContract.Address, nil)
	if err != nil {
		return 0, errors.Wrapf(err, "Fail to get liqBalance. operatorId:%d", operatorId)
	}

	operatorAvailable := big.NewInt(0).Add(operatorPoolBalances, operatorNftPoolBalances)
	// If the contract funds are borrowed, it is impossible to support the operator available so many deposits,
	// and the maximum contract balance can only be used
	if liqBalance.Cmp(operatorAvailable) == -1 {
		operatorAvailable = liqBalance
	}
	count := new(big.Int).Div(operatorAvailable, eth1.ETH32())

	return count.Uint64(), nil
}

// RegisterValidatorForNodeDAO register validator for node dao
// @return [][]byte pubkeys
func (e *RegisterValidator) RegisterValidatorForNodeDAO(operatorId int64, depositDatas []*DepositData) ([][]byte, error) {

	var pubkeys = make([][]byte, 0)
	var signatures = make([][]byte, 0)
	var depositDataRoots = make([][32]byte, 0)

	for _, depositData := range depositDatas {
		pubkey, err := hex.DecodeString(strings.TrimPrefix(depositData.Pubkey, "0x"))
		if err != nil {
			return nil, errors.Wrapf(err, "Fail to decode pubkey. pubkey:%v", depositData.Pubkey)
		}
		pubkeys = append(pubkeys, pubkey)

		signture, err := hex.DecodeString(strings.TrimPrefix(depositData.Signature, "0x"))
		if err != nil {
			return nil, errors.Wrapf(err, "Fail to decode signture. signture:%v", depositData.Signature)
		}
		signatures = append(signatures, signture)

		depositDataRoot, err := hex.DecodeString(strings.TrimPrefix(depositData.DepositDataRoot, "0x"))
		if err != nil {
			return nil, errors.Wrapf(err, "Fail to decode depositDataRoot. signture:%v", depositData.DepositDataRoot)
		}

		var root [32]byte
		copy(root[:], depositDataRoot[:32])
		depositDataRoots = append(depositDataRoots, root)
	}

	tx, err := e.liqContract.Contract.RegisterValidator(e.keyTransactOpts, pubkeys, signatures, depositDataRoots)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to register validator. operatorId:%d", operatorId)
	}

	// Wait for the transaction to complete
	receipt, err := bind.WaitMined(context.Background(), e.elClient.Client, tx)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to WaitMined RegisterValidator. tx hash:%s, operatorId:%v", tx.Hash().String(), operatorId)
	}

	isSuccess := eth1.TxReceiptSuccess(receipt)
	if !isSuccess {
		return nil, errors.Errorf("RegisterValidator tx failed. tx hash:%s, operatorId:%v", tx.Hash().String(), operatorId)
	}

	logger.Infof("RegisterValidator success. tx hash:%s, operatorId:%v", tx.Hash().String(), operatorId)

	return pubkeys, nil
}
