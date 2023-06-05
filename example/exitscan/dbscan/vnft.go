// description:
// @author renshiwei
// Date: 2023/6/5

package dbscan

import (
	"context"
	"github.com/NodeDAO/operator-sdk/validator/exitscan"
	"github.com/pkg/errors"
	"math/big"
)

type DBVnftExitter struct {
	exitScanner exitscan.ExitScanner
	network     string
}

func NewDBVnftExitFilter(ctx context.Context, network, elAddr string) (*DBVnftExitter, error) {
	vnftExitScan, err := exitscan.NewVnftExitScan(ctx, network, elAddr)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to new vnftExitScan")
	}

	return &DBVnftExitter{
		exitScanner: vnftExitScan,
		network:     network,
	}, nil
}

func (e *DBVnftExitter) Filter(operatorId *big.Int) ([]interface{}, error) {
	vnftContractExitRecords, err := e.exitScanner.ExitScan(operatorId)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to ExitScan. network:%s operatorId:%s", e.network, operatorId.String())
	}

	tokenIds := make([]*big.Int, 0, len(vnftContractExitRecords))

	for _, record := range vnftContractExitRecords {
		tokenIds = append(tokenIds, record.TokenId)
	}

	recordDao := NewVnftRecord()
	nodedaoValidators, err = recordDao.GetByTokenIds(e.network, tokenIds, false)
	if err != nil {
		return nil, errors.Wrapf(err, "Fail to GetByTokenIds by db. network:%s operatorId:%s", e.network, operatorId.String())
	}

	validators := make([]interface{}, len(vnftContractExitRecords))
	for i, record := range nodedaoValidators {
		validators[i] = record
	}

	return validators, nil
}

func (e *DBVnftExitter) ExitMark(operatorId *big.Int, vnftRecords []interface{}) error {

	return nil
}
