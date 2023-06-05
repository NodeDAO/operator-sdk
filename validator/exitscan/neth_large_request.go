// description:
// @author renshiwei
// Date: 2023/5/24

package exitscan

import (
	"context"
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

// WithdrawalRequestScan Scanning the smart contract requires processing the exiting WithdrawalRequest
// !!! Handling exits is delayed, and additional operations are required to mark and filter the WithdrawalRequest,
// the simplest way is to use db, see example for this part
func (s *NETHExitScan) WithdrawalRequestScan(operatorId *big.Int) ([]*withdrawalRequest.WithdrawalRequestWithdrawalInfo, error) {
	withdrawalInfos := make([]*withdrawalRequest.WithdrawalRequestWithdrawalInfo, 0)

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
			withdrawalInfos = append(withdrawalInfos, &withdrawalRequest.WithdrawalRequestWithdrawalInfo{
				OperatorId:         opId,
				WithdrawHeight:     withdrawHeight,
				WithdrawNethAmount: withdrawNethAmount,
				WithdrawExchange:   withdrawExchange,
				ClaimEthAmount:     claimEthAmount,
				Owner:              owner,
				IsClaim:            isClaim,
			})
		}

		i++
	}

	return withdrawalInfos, nil
}

// ExitScan Filter for exits
// !!! Use the filtered []WithdrawalRequest to operate
// @param operatorId operator id
func (s *NETHExitScan) ExitScan(operatorId *big.Int) ([]*VnftRecord, error) {
	vnfts := make([]*VnftRecord, 0)

	// Get the number of active vnft of the operator.
	vnftActiveCount, err := s.vnftContract.Contract.GetUserActiveNftCountsOfOperator(nil, operatorId)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get vnftContract GetUserActiveNftCountsOfOperator. network:%s", s.network)
	}
	if vnftActiveCount.Cmp(big.NewInt(0)) == 0 {
		return vnfts, nil
	}

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
				OperatorId: operatorId,
				TokenId:    tokenId,
				Pubkey:     pubkey,
				Type:       LiquidStaking,
			})
		}

		if vnftActiveCount.Cmp(big.NewInt(int64(len(vnfts)))) == 0 {
			break
		}
	}

	return vnfts, nil
}
