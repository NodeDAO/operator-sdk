// description:
// @author renshiwei
// Date: 2023/6/7

package dbscan

import (
	"github.com/NodeDAO/operator-sdk/common/logger"
	"github.com/NodeDAO/operator-sdk/example/dao"
	"github.com/NodeDAO/operator-sdk/validator/exitscan"
	"github.com/pkg/errors"
	"math/big"
)

type DBMarker interface {
	exitscan.ExitMarker
	exitscan.WithdrawalRequestMarker
}

type DBMark struct {
	network string
}

func NewDBMark(network string) (*DBMark, error) {
	return &DBMark{
		network: network,
	}, nil
}

// ExitMark Mark the filtered Vnft Records as exiting for the next filter
// @param vnftRecords VNFT Records that need to be tagged
func (e *DBMark) ExitMark(operatorId *big.Int, vnftRecords []*exitscan.VnftRecord) error {
	tokenIds := make([]uint64, 0, len(vnftRecords))
	for _, record := range vnftRecords {
		tokenIds = append(tokenIds, record.TokenId.Uint64())
	}

	recordDao := dao.NewNodedaoValidator()
	err := recordDao.UpdateExited(e.network, tokenIds)
	if err != nil {
		return errors.Wrapf(err, "Fail to UpdateExited by db. network:%s operatorId:%s", e.network, operatorId.String())
	}

	logger.Infof("ExitMark success by DB. tokenIds:%+v", tokenIds)

	return nil
}

// WithdrawalRequestMark Mark the filtered withdrawalRequests as deal for the next filter
// @param withdrawalRequests that need to be tagged
func (e *DBMark) WithdrawalRequestMark(operatorId *big.Int, withdrawalRequests []*exitscan.WithdrawalRequest) error {
	withdrawalRequestIds := make([]*big.Int, 0, len(withdrawalRequests))
	for _, record := range withdrawalRequests {
		withdrawalRequestIds = append(withdrawalRequestIds, record.ID)
	}

	withdrawalRequestDao := dao.NewNethWithdrawalRequest()
	err := withdrawalRequestDao.UpdateExited(e.network, withdrawalRequestIds)
	if err != nil {
		return errors.Wrapf(err, "Fail to UpdateExited by db. network:%s operatorId:%s", e.network, operatorId.String())
	}

	logger.Infof("WithdrawalRequestMark success by DB. withdrawalRequestIds:%+v", withdrawalRequestIds)

	return nil
}
