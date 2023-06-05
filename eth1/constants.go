// desc:
// @author renshiwei
// Date: 2023/4/7 19:31

package eth1

import (
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

const (
	MAINNET = "mainnet"
	GOERLI  = "goerli"
	PRATER  = "prater"
	SEPOLIA = "sepolia"
)

const (
	// CONTRACT_REVENT The contract revent reported an incorrect keyword
	CONTRACT_REVENT = "execution reverted"
)

var ZERO_HASH = [32]byte{}

const ZERO_HASH_STR = "0x0000000000000000000000000000000000000000000000000000000000000000"

func GWEIToWEI(value *big.Int) *big.Int {
	return new(big.Int).Mul(value, big.NewInt(params.GWei))
}

func ETH32() *big.Int {
	eth32, _ := new(big.Int).SetString("32000000000000000000", 0)
	return eth32
}
