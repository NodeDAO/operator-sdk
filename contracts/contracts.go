// desc:
// @author renshiwei
// Date: 2023/4/7 11:38

package contracts

import (
	"github.com/NodeDAO/operator-sdk/contracts/liq"
	"github.com/NodeDAO/operator-sdk/contracts/vnft"
	"github.com/NodeDAO/operator-sdk/contracts/withdrawalRequest"
	"github.com/NodeDAO/operator-sdk/eth1"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"strings"
)

type VnftContract struct {
	Network  string
	Address  string
	Contract *vnft.Vnft
}

type LiqContract struct {
	Network  string
	Address  string
	Contract *liq.Liq
}

type WithdrawalRequestContract struct {
	Network  string
	Address  string
	Contract *withdrawalRequest.WithdrawalRequest
}

const (
	LIQ_ADDRESS_MAINNET = "0x8103151E2377e78C04a3d2564e20542680ed3096"
	LIQ_ADDRESS_GOERLI  = "0x949AC43bb71F8710B0F1193880b338f0323DeB1a"

	VNFT_ADDRESS_MAINNET = "0x58553F5c5a6AEE89EaBFd42c231A18aB0872700d"
	VNFT_ADDRESS_GOERLI  = "0x3CB42bb75Cf1BcC077010ac1E3d3Be22D13326FA"

	WITHDRAWAL_REQUEST_ADDRESS_MAINNET = "0xE81fC969D14Cad8537ebAFa2a1c478F29d7840FC"
	WITHDRAWAL_REQUEST_ADDRESS_GOERLi  = "0x006e69F509E31c91263C03a744B47c3b03eAC391"
)

func NewVnftContract(network string, eth1Client bind.ContractBackend) (*VnftContract, error) {
	e := &VnftContract{
		Network: network,
	}
	if strings.ToLower(network) == eth1.MAINNET {
		e.Address = VNFT_ADDRESS_MAINNET
	} else if strings.ToLower(network) == eth1.GOERLI {
		e.Address = VNFT_ADDRESS_GOERLI
	}

	if e.Address == "" {
		return nil, errors.New("vnftContract contract address is empty.")
	}

	var err error
	e.Contract, err = vnft.NewVnft(common.HexToAddress(e.Address), eth1Client)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to new VNFT.")
	}
	return e, nil
}

func NewLiqContract(network string, eth1Client bind.ContractBackend) (*LiqContract, error) {
	e := &LiqContract{
		Network: network,
	}
	if strings.ToLower(network) == eth1.MAINNET {
		e.Address = LIQ_ADDRESS_MAINNET
	} else if strings.ToLower(network) == eth1.GOERLI {
		e.Address = LIQ_ADDRESS_GOERLI
	}

	if e.Address == "" {
		return nil, errors.New("liqContract contract address is empty.")
	}

	var err error
	e.Contract, err = liq.NewLiq(common.HexToAddress(e.Address), eth1Client)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to new liqContract.")
	}
	return e, nil
}

func NewWithdrawalRequestContract(network string, eth1Client bind.ContractBackend) (*WithdrawalRequestContract, error) {
	e := &WithdrawalRequestContract{
		Network: network,
	}
	if strings.ToLower(network) == eth1.MAINNET {
		e.Address = WITHDRAWAL_REQUEST_ADDRESS_MAINNET
	} else if strings.ToLower(network) == eth1.GOERLI {
		e.Address = WITHDRAWAL_REQUEST_ADDRESS_GOERLi
	}

	var err error
	e.Contract, err = withdrawalRequest.NewWithdrawalRequest(common.HexToAddress(e.Address), eth1Client)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to new withdrawal Request.")
	}
	return e, nil
}
