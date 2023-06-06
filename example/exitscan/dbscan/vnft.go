// description:
// @author renshiwei
// Date: 2023/6/5

package dbscan

import (
	"github.com/NodeDAO/operator-sdk/example/dao"
	"github.com/NodeDAO/operator-sdk/validator/exitscan"
	"github.com/pkg/errors"
	"math/big"
)

type DBVnftExitter struct {
	network string
}

func NewDBVnftExitFilter(network string) (*DBVnftExitter, error) {
	return &DBVnftExitter{
		network: network,
	}, nil
}

func (e *DBVnftExitter) Filter(operatorId *big.Int, vnftContractExitRecords []*exitscan.VnftRecord) ([]*exitscan.VnftRecord, error) {
	tokenIds := make([]*big.Int, 0, len(vnftContractExitRecords))
	for _, record := range vnftContractExitRecords {
		tokenIds = append(tokenIds, record.TokenId)
	}

	recordDao := dao.NewNodedaoValidator()
	nodedaoValidators, err := recordDao.GetByTokenIds(e.network, tokenIds, false)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to GetByTokenIds by db. network:%s operatorId:%s", e.network, operatorId.String())
	}

	validators := make([]*exitscan.VnftRecord, len(nodedaoValidators))
	for i, record := range nodedaoValidators {
		validators[i] = &exitscan.VnftRecord{
			Network:    e.network,
			OperatorId: operatorId,
			TokenId:    record.TokenId,
			Pubkey:     record.Pubkey,
			Type:       record.Type,
		}
	}
	return validators, nil
}

func (e *DBVnftExitter) ExitMark(operatorId *big.Int, vnftRecords []*exitscan.VnftRecord) error {
	tokenIds := make([]*big.Int, 0, len(vnftRecords))
	for _, record := range vnftRecords {
		tokenIds = append(tokenIds, record.TokenId)
	}

	recordDao := dao.NewNodedaoValidator()
	err := recordDao.UpdateExited(e.network, tokenIds)
	if err != nil {
		return errors.Wrapf(err, "Fail to UpdateExited by db. network:%s operatorId:%s", e.network, operatorId.String())
	}

	return nil
}
