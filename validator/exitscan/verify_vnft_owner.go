// description:
// @author renshiwei
// Date: 2023/6/7

package exitscan

import (
	"context"
	"github.com/NodeDAO/operator-sdk/contracts"
	"github.com/NodeDAO/operator-sdk/eth1"
	"github.com/pkg/errors"
	"math/big"
)

type VnftOwnerVerify struct {
	// param
	network string

	// init
	elClient *eth1.EthClient
	// contracts
	vnftContract *contracts.VnftContract
}

// NewVnftOwnerVerify new vnft owner verify
func NewVnftOwnerVerify(ctx context.Context, network, elAddr string) (*VnftOwnerVerify, error) {
	var err error

	elClient, err := eth1.NewEthClient(ctx, elAddr)
	if err != nil {
		return nil, err
	}

	vnftContract, err := contracts.NewVnftContract(network, elClient.Client)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to new vnftContract. network:%s", network)
	}

	return &VnftOwnerVerify{
		network:      network,
		elClient:     elClient,
		vnftContract: vnftContract,
	}, nil
}

// VerifyStakeType Verify that the stakeType of vnft tokenIds
func (e *VnftOwnerVerify) VerifyStakeType(network string, tokenId *big.Int) (StakeType, error) {
	nETHOwnerExpected := contracts.LiqContractAddress(network)

	vnftOwner, err := e.vnftContract.Contract.OwnerOf(nil, tokenId)
	if err != nil {
		return 0, errors.Wrapf(err, "Fail to get vnft owner. network:%s", e.network)
	}

	if vnftOwner.String() != nETHOwnerExpected {
		return NETH, nil
	}

	return VNFT, nil
}

// VerifyVnftOwner Verify that the stakeType of vnft tokenIds and vnftOwner match
// ----------------------------------------------------------------
// The relationship between StakeType and VnftOwner is as follows:
// ----------------------
// StakeType | VnftOwner
// ----------------------
// VNFT      | USER
// NETH      | LiquidStaking
func (e *VnftOwnerVerify) VerifyVnftOwner(network string, stakeType StakeType, vnftOwner VnftOwner, tokenIds []*big.Int) (bool, error) {
	if uint32(stakeType) != uint32(vnftOwner) {
		return false, errors.New("stakeType and vnftOwner do not match")
	}

	var nETHOwnerExpected string
	if stakeType == NETH {
		nETHOwnerExpected = contracts.LiqContractAddress(network)
	}

	for _, tokenId := range tokenIds {
		vnftOwner, err := e.vnftContract.Contract.OwnerOf(nil, tokenId)
		if err != nil {
			return false, errors.Wrapf(err, "Fail to get vnft owner. network:%s", e.network)
		}

		if stakeType == NETH {
			if vnftOwner.String() != nETHOwnerExpected {
				return false, errors.New("stakeType == NETH, vnft owner is not LiquidStaking")
			}
		} else if stakeType == VNFT {
			if vnftOwner.String() == nETHOwnerExpected {
				return false, errors.New("stakeType == VNFT, vnft owner is LiquidStaking")
			}
		}

	}

	return true, nil
}
