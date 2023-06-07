// description:
// @author renshiwei
// Date: 2023/5/24

// Package exitscan NodeDAO operator exit scan
package exitscan

import (
	"fmt"
	"github.com/NodeDAO/operator-sdk/contracts/withdrawalRequest"
	"math/big"
)

type VnftOwner uint32
type StakeType uint32

const (
	USER VnftOwner = iota
	LiquidStaking
)

const (
	VNFT StakeType = iota
	NETH
)

func GetVnftOwner(stakeType StakeType) VnftOwner {
	if stakeType == VNFT {
		return USER
	} else {
		return LiquidStaking
	}
}

type VnftRecord struct {
	Network    string
	OperatorId *big.Int
	TokenId    *big.Int
	Pubkey     string
	Type       StakeType
}

func (v *VnftRecord) String() string {
	return fmt.Sprintf("{Network:%s OperatorId:%s, TokenId:%s, Pubkey:%s, Type:%v}", v.Network, v.OperatorId.String(), v.TokenId.String(), v.Pubkey, v.Type)
}

type WithdrawalRequest struct {
	ID                    *big.Int
	WithdrawalRequestInfo *withdrawalRequest.WithdrawalRequestWithdrawalInfo
}

type VnftOwnerValidator interface {
	// VerifyVnftOwner Verify that the stakeType of vnft tokenIds and vnftOwner match
	// ----------------------------------------------------------------
	// The relationship between StakeType and VnftOwner is as follows:
	// ----------------------
	// StakeType | VnftOwner
	// ----------------------
	// VNFT      | USER
	// NETH      | LiquidStaking
	VerifyVnftOwner(network string, stakeType StakeType, vnftOwner VnftOwner, tokenIds []*big.Int) (bool, error)
}

// ExitScanner Scan the smart contract for records that need to be exited
type ExitScanner interface {
	ExitScan(operatorId *big.Int) ([]*VnftRecord, error)
}

// WithdrawalRequestScanner nETH exit depends on the WithdrawalRequest.
// vNFT exit can be used directly with Exit Scanner
type WithdrawalRequestScanner interface {
	ExitScanner

	// WithdrawalRequestScan Scan for unclaimed Withdrawal Requests
	WithdrawalRequestScan(operatorId *big.Int) ([]*withdrawalRequest.WithdrawalRequestWithdrawalInfo, error)
}

// ExitFilter filter the exit for vNFT and nETH.
// Validator's exit is asynchrony. The reasons for asynchrony are:
// 1. The validator exit goes through the lifetime of the beacon
// 2. NodeDAO-Oracle is required to report to settle
// --------------------------------------------
// Filter The operator needs to implement it by itself, and the easiest way is to use db.
// An example implementation will be provided, based on MySQL, see Example
type ExitFilter interface {
	// Filter To filter for exit
	// @return []*VnftRecord{} Filtered
	Filter(operatorId *big.Int, vnftRecords []*VnftRecord) ([]*VnftRecord, error)
}

// WithdrawalRequestFilter To filter for WithdrawalRequests
// --------------------------------------------
// The simplest way to implement the operator is to use db, see example
type WithdrawalRequestFilter interface {
	// WithdrawalRequestFilter To filter for WithdrawalRequests
	// @return []*WithdrawalRequest Filtered WithdrawalRequests
	WithdrawalRequestFilter(operatorId *big.Int, withdrawalRequests []*WithdrawalRequest) ([]*WithdrawalRequest, error)
}

// ExitMarker To perform a validator exit, it needs to be flagged, and then it is used for filter
// --------------------------------------------
// The simplest way to implement the operator is to use db, see example
type ExitMarker interface {
	// ExitMark Mark the exit of the Vnft Record
	ExitMark(operatorId *big.Int, vnftRecords []*VnftRecord) error
}

// WithdrawalRequestMarker To mark deal for WithdrawalRequest
// --------------------------------------------
// The simplest way to implement the operator is to use db, see example
type WithdrawalRequestMarker interface {
	// WithdrawalRequestMark To mark deal for WithdrawalRequest
	WithdrawalRequestMark(operatorId *big.Int, withdrawalRequests []*WithdrawalRequest) error
}
