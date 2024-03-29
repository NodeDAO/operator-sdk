// description:
// @author renshiwei
// Date: 2023/5/24

package exitscan

import (
	"context"
	"github.com/NodeDAO/operator-sdk/common/logger"
	"github.com/NodeDAO/operator-sdk/contracts"
	"github.com/NodeDAO/operator-sdk/contracts/withdrawalRequest"
	"github.com/NodeDAO/operator-sdk/eth1"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"math/big"
	"strings"
)

// NETHExitScan nETH exit scan
// Only largeRequest nETH needs to be scanned to exit
type NETHExitScan struct {
	// param
	network string

	// init
	elClient *eth1.EthClient
	// contracts
	withdrawalRequestContract *contracts.WithdrawalRequestContract
	vnftContract              *contracts.VnftContract
	liqContract               *contracts.LiqContract
}

// NewNETHExitScan new nETH exit scan
func NewNETHExitScan(ctx context.Context, network, elAddr string) (*NETHExitScan, error) {
	var err error

	elClient, err := eth1.NewEthClient(ctx, elAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to new eth client. network:%s", network)
	}

	withdrawalRequestContract, err := contracts.NewWithdrawalRequestContract(network, elClient.Client)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to new withdrawalRequestContract. network:%s", network)
	}
	vnftContract, err := contracts.NewVnftContract(network, elClient.Client)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to new vnftContract. network:%s", network)
	}
	liqContract, err := contracts.NewLiqContract(network, elClient.Client)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to new liqContract. network:%s", network)
	}

	return &NETHExitScan{
		network:                   network,
		elClient:                  elClient,
		withdrawalRequestContract: withdrawalRequestContract,
		vnftContract:              vnftContract,
		liqContract:               liqContract,
	}, nil
}

// WithdrawalRequestScan Scanning the smart contract requires processing the exiting WithdrawalRequest.
// !!! Handling exits is delayed, and additional operations are required to mark and filter the WithdrawalRequest,
// the simplest way is to use db, see example for this part.
func (s *NETHExitScan) WithdrawalRequestScan(operatorId *big.Int) ([]*WithdrawalRequest, error) {
	withdrawalInfos := make([]*WithdrawalRequest, 0)

	// Iterate through the withdrawal Request in the contract that is not handled by the current operator.
	// The current withdrawalRequestContract method 'getWithdrawalOfOperator' does not return a requestId and cannot be used at the moment.
	i := 0
	for {
		opId, withdrawHeight, withdrawNethAmount, withdrawExchange, claimEthAmount, owner, isClaim, err := s.withdrawalRequestContract.Contract.GetWithdrawalOfRequestId(nil, big.NewInt(int64(i)))

		if err != nil {
			if strings.Contains(err.Error(), eth1.CONTRACT_REVENT) {
				break
			} else {
				return nil, errors.Wrapf(err, "Fail to get withdrawalRequestContract GetWithdrawalOfRequestId. network:%s", s.network)
			}
		}

		// If the operatorId is the same as the current operatorId, and the withdrawal has not been claimed, it needs to be exited\
		if operatorId.Cmp(opId) == 0 && !isClaim {
			withdrawalInfos = append(withdrawalInfos, &WithdrawalRequest{
				ID: big.NewInt(int64(i)),
				WithdrawalRequestInfo: &withdrawalRequest.WithdrawalRequestWithdrawalInfo{
					OperatorId:         opId,
					WithdrawHeight:     withdrawHeight,
					WithdrawNethAmount: withdrawNethAmount,
					WithdrawExchange:   withdrawExchange,
					ClaimEthAmount:     claimEthAmount,
					Owner:              owner,
					IsClaim:            isClaim,
				},
			})
		}

		i++
	}

	logger.Infof("nETH ExitScan success by contract. withdrawalInfos:%+v", withdrawalInfos)

	return withdrawalInfos, nil
}

// ExitScan Filter for exits
// !!! Use the filtered []WithdrawalRequest to operator
// @param operatorId operator id
func (s *NETHExitScan) ExitScan(operatorId *big.Int) ([]*VnftRecord, error) {
	vnfts := make([]*VnftRecord, 0)

	// Query the tokenId of all vnfts owned by the LiquidStaking pool
	stakingPoolTokenIds, err := s.vnftContract.Contract.ActiveNftsOfStakingPool(nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get vnftContract ActiveNftsOfStakingPool. network:%s", s.network)
	}

	for _, tokenId := range stakingPoolTokenIds {
		// Filter based on operatorId
		operatorOf, err := s.vnftContract.Contract.OperatorOf(nil, tokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "Fail to get vnftContract OperatorOf. network:%s", s.network)
		}
		// If the operatorId is the same as the current operatorId, and the withdrawal has not been claimed, it needs to be exited
		if operatorOf.Cmp(operatorId) == 0 {
			pubkeyBytes, err := s.vnftContract.Contract.ValidatorOf(nil, tokenId)
			if err != nil {
				return nil, errors.Wrapf(err, "Fail to get vnftContract ValidatorOf. network:%s", s.network)
			}
			pubkey := hexutil.Encode(pubkeyBytes)

			vnfts = append(vnfts, &VnftRecord{
				Network:    s.network,
				OperatorId: operatorId,
				TokenId:    tokenId,
				Pubkey:     pubkey,
				Type:       NETH,
			})
		}

	}

	logger.Infof("nETH ExitScan success by contract. stakingPoolTokenIds:%+v", stakingPoolTokenIds)

	return vnfts, nil
}

// ExitCounter Calculate the number of validators that need to be exited by a Withdrawal Request
// @param filterWithdrawalRequests  A list of offline filtered Withdrawal Requests
// --------------------------------------------------------
// if sumETHAmount = 64 ether, need to exit 2 validator
// if sumETHAmount = 66 ether, need to exit 3 validator
func (s *NETHExitScan) ExitCounter(filterWithdrawalRequests []*WithdrawalRequest) (uint32, error) {
	if len(filterWithdrawalRequests) == 0 {
		return 0, nil
	}

	sumETHAmount := big.NewInt(0)
	for _, request := range filterWithdrawalRequests {
		sumETHAmount = new(big.Int).Add(sumETHAmount, request.WithdrawalRequestInfo.ClaimEthAmount)
	}

	if sumETHAmount.Cmp(big.NewInt(0)) == 0 {
		return 0, nil
	}

	// Calculate the number of withdrawals The part greater than 32eth needs to withdraw one more
	vnftScanCount := 0
	div, mod := new(big.Int).DivMod(sumETHAmount, eth1.ETH32(), new(big.Int))
	if mod.Cmp(big.NewInt(0)) == 1 {
		vnftScanCount = int(div.Uint64()) + 1
	} else {
		vnftScanCount = int(div.Uint64())
	}

	return uint32(vnftScanCount), nil
}
