// description:
// @author renshiwei
// Date: 2023/6/6

package dbscan

import "github.com/NodeDAO/operator-sdk/validator/exitscan"

type DBNethExitter struct {
	network string
}

func NewDBNethExitFilter(network string) (*DBNethExitter, error) {
	return &DBNethExitter{
		network: network,
	}, nil
}

func (e *DBNethExitter) Filter(operatorId string, nethContractExitRecords []*exitscan.VnftRecord) ([]*exitscan.VnftRecord, error) {

	return nil, nil
}
