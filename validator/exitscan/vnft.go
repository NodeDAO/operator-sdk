// description:
// @author renshiwei
// Date: 2023/5/24

package exitscan

import (
	"context"
	"github.com/NodeDAO/operator-sdk/contracts"
	"github.com/NodeDAO/operator-sdk/eth1"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"math/big"
)

// VnftExitScan vnft exit scan
type VnftExitScan struct {
	// param
	network string

	// init
	elClient *eth1.EthClient
	// contracts
	withdrawalRequestContract *contracts.WithdrawalRequestContract
	vnftContract              *contracts.VnftContract
}

// NewVnftExitScan new vnft exit scan
func NewVnftExitScan(ctx context.Context, network, elAddr string) (*VnftExitScan, error) {
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

	return &VnftExitScan{
		network:                   network,
		elClient:                  elClient,
		withdrawalRequestContract: withdrawalRequestContract,
		vnftContract:              vnftContract,
	}, nil
}

// ExitScan get contract need exit vnft list
func (s *VnftExitScan) ExitScan(operatorId *big.Int) ([]*VnftRecord, error) {
	tokenIds, err := s.withdrawalRequestContract.Contract.GetUserUnstakeButOperatorNoExitNfs(nil, operatorId)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to get withdrawalRequestContract GetUserUnstakeButOperatorNoExitNfs. network:%s", s.network)
	}

	vnfts := make([]*VnftRecord, 0)

	for _, tokenId := range tokenIds {
		pubkeyBytes, err := s.vnftContract.Contract.ValidatorOf(nil, tokenId)
		if err != nil {
			return nil, errors.Wrapf(err, "Fail to get vnftContract ValidatorOf. network:%s", s.network)
		}
		pubkey := hexutil.Encode(pubkeyBytes)

		vnft := &VnftRecord{
			Network:    s.network,
			OperatorId: operatorId,
			TokenId:    tokenId,
			Pubkey:     pubkey,
			Type:       USER,
		}

		vnfts = append(vnfts, vnft)
	}

	return vnfts, nil
}
