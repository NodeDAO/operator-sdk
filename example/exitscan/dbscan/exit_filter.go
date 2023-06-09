// description:
// @author renshiwei
// Date: 2023/6/7

package dbscan

import (
	"github.com/NodeDAO/operator-sdk/common/logger"
	"github.com/NodeDAO/operator-sdk/contracts/withdrawalRequest"
	"github.com/NodeDAO/operator-sdk/example/dao"
	"github.com/NodeDAO/operator-sdk/validator/exitscan"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

type DBFilter interface {
	exitscan.ExitFilter
	exitscan.WithdrawalRequestFilter
}

type DBFiltrate struct {
	network   string
	stakeType exitscan.StakeType
	vnftOwner exitscan.VnftOwner

	vnftOwnerValidator       exitscan.VnftOwnerValidator
	nethExitValidatorCounter exitscan.WithdrawalRequestExitValidatorCounter

	// mid res
	exitValidatorCount uint32
}

func NewDBFiltrate(
	network string,
	stakeType exitscan.StakeType,
	vnftOwnerValidator exitscan.VnftOwnerValidator,
) (*DBFiltrate, error) {
	return &DBFiltrate{
		network:            network,
		stakeType:          stakeType,
		vnftOwner:          exitscan.GetVnftOwner(stakeType),
		vnftOwnerValidator: vnftOwnerValidator,
	}, nil
}

func (e *DBFiltrate) SetExitValidatorCounter(nethExitValidatorCounter exitscan.WithdrawalRequestExitValidatorCounter) {
	e.nethExitValidatorCounter = nethExitValidatorCounter
}

// WithdrawalRequestFilter 1. Filter out the withdrawalRequests that have been exited.
// 2. and calculate the number of validators that nETH needs to exit, according to param withdrawalRequests
// ------------------------------------------------------------------------------------------------
// !!! Call 'WithdrawalRequestFilter' before calling 'SetExitValidatorCounter'
func (e *DBFiltrate) WithdrawalRequestFilter(operatorId *big.Int, withdrawalRequests []*exitscan.WithdrawalRequest) ([]*exitscan.WithdrawalRequest, error) {
	withdrawalRequestDao := dao.NewNethWithdrawalRequest()

	withdrawalRequestIds := make([]uint64, 0, len(withdrawalRequests))
	saveWithdrawals := make([]*dao.NethWithdrawalRequest, 0, len(withdrawalRequests))
	for _, record := range withdrawalRequests {
		withdrawalRequestIds = append(withdrawalRequestIds, record.ID.Uint64())
		saveWithdrawals = append(saveWithdrawals, &dao.NethWithdrawalRequest{
			Network:            e.network,
			OperatorId:         record.WithdrawalRequestInfo.OperatorId.Uint64(),
			RequestId:          record.ID.Uint64(),
			WithdrawHeight:     record.WithdrawalRequestInfo.WithdrawHeight.String(),
			WithdrawNethAmount: record.WithdrawalRequestInfo.WithdrawNethAmount.String(),
			WithdrawExchange:   record.WithdrawalRequestInfo.WithdrawExchange.String(),
			ClaimEthAmount:     record.WithdrawalRequestInfo.ClaimEthAmount.String(),
			Owner:              record.WithdrawalRequestInfo.Owner.String(),
		})
	}

	// largeRequest if it does not exist, it is inserted
	err := withdrawalRequestDao.InsertForNotExist(saveWithdrawals)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to InsertForNotExist by db. network:%s operatorId:%s", e.network, operatorId.String())
	}

	// The odds of withdrawalRequests are only queried in the database for is_exit = false.
	exitWithdrawalRequests, err := withdrawalRequestDao.GetByRequestIds(e.network, operatorId.Uint64(), withdrawalRequestIds, false)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to GetByRequestIds by db. network:%s operatorId:%s", e.network, operatorId.String())
	}

	withdrawals := make([]*exitscan.WithdrawalRequest, len(exitWithdrawalRequests))
	for i, record := range exitWithdrawalRequests {

		withdrawNethAmount, _ := new(big.Int).SetString(record.WithdrawNethAmount, 0)
		withdrawExchange, _ := new(big.Int).SetString(record.WithdrawExchange, 0)
		claimEthAmount, ok := new(big.Int).SetString(record.ClaimEthAmount, 0)
		withdrawHeight, _ := new(big.Int).SetString(record.WithdrawHeight, 0)
		if !ok {
			return nil, errors.New("Fail to string cast big.Int.")
		}

		withdrawals[i] = &exitscan.WithdrawalRequest{
			ID: big.NewInt(int64(record.RequestId)),
			WithdrawalRequestInfo: &withdrawalRequest.WithdrawalRequestWithdrawalInfo{
				OperatorId:         big.NewInt(int64(record.OperatorId)),
				WithdrawHeight:     withdrawHeight,
				WithdrawNethAmount: withdrawNethAmount,
				WithdrawExchange:   withdrawExchange,
				ClaimEthAmount:     claimEthAmount,
				Owner:              common.HexToAddress(record.Owner),
			},
		}
	}

	counter, err := e.nethExitValidatorCounter.ExitCounter(withdrawals)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to ExitCounter by nethExitValidatorCounter.")
	}
	e.exitValidatorCount = counter

	logger.Infof("Filter success by DB. ExitCounter:%v withdrawals:%+v", counter, withdrawals)

	return withdrawals, nil
}

// Filter If some validator records for vnft ContractExitRecords have been marked as is_exit = true in the database, filtering.
// @param vnftContractExitRecords *[]exitscan.VnftRecord Unfiltered VnftRecord.
// @return []*exitscan.VnftRecord Filtered VnftRecord.
// ----------------------------------------------------------------
// if dbCanExitValidatorCount!=needExitValidatorCount, return error
// The exitValidatorCount calculated by 'WithdrawalRequestFilter' filters out the specified number of VnftRecords again
func (e *DBFiltrate) Filter(operatorId *big.Int, vnftContractExitRecords []*exitscan.VnftRecord) ([]*exitscan.VnftRecord, error) {
	tokenIds := make([]uint64, 0, len(vnftContractExitRecords))
	tokenIdBigInts := make([]*big.Int, 0, len(vnftContractExitRecords))
	for _, record := range vnftContractExitRecords {
		tokenIds = append(tokenIds, record.TokenId.Uint64())
		tokenIdBigInts = append(tokenIdBigInts, record.TokenId)
	}

	recordDao := dao.NewNodedaoValidator()

	// The odds of vnftContractExitRecords are only queried in the database for is_exit = false.
	nodedaoValidators, err := recordDao.GetByTokenIds(e.network, operatorId.Uint64(), tokenIds, e.stakeType, false)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to GetByTokenIds by db. network:%s operatorId:%s", e.network, operatorId.String())
	}

	// Verify the ownership of tokenId again
	isVerify, err := e.vnftOwnerValidator.VerifyVnftOwner(e.network, e.stakeType, e.vnftOwner, tokenIdBigInts)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to VerifyVnftOwner by db. network:%s operatorId:%s", e.network, operatorId.String())
	}
	if !isVerify {
		return nil, errors.Errorf("Fail to VerifyVnftOwner.stakeType:%v tokenIds:%+v", e.stakeType, tokenIds)
	}

	exitValidatorCount := 0
	if e.stakeType == exitscan.NETH {
		if len(nodedaoValidators) < int(e.exitValidatorCount) {
			return nil, errors.Errorf("Fail to Filter. dbCanExitValidatorCount!=needExitValidatorCount. dbCanExitValidatorCount:%v needExitValidatorCount:%v", len(nodedaoValidators), e.exitValidatorCount)
		}
		exitValidatorCount = int(e.exitValidatorCount)
	} else if e.stakeType == exitscan.VNFT {
		exitValidatorCount = len(nodedaoValidators)
	}

	validators := make([]*exitscan.VnftRecord, exitValidatorCount)
	for i, record := range nodedaoValidators {
		validators[i] = &exitscan.VnftRecord{
			Network:    e.network,
			OperatorId: operatorId,
			TokenId:    big.NewInt(int64(record.TokenId)),
			Pubkey:     record.Pubkey,
			Type:       record.Type,
		}
	}

	logger.Infof("Filter success by DB. validators:%+v", validators)

	return validators, nil
}
