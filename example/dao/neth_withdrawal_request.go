// description:
// @author renshiwei
// Date: 2023/6/6

package dao

import (
	"github.com/NodeDAO/operator-sdk/config"
	"github.com/pkg/errors"
	"math/big"
	"time"
)

type NethWithdrawalRequest struct {
	ID                 uint64    `gorm:"column:id;primary_key" json:"id"`
	Network            string    `gorm:"column:network;" json:"network"`
	OperatorId         *big.Int  `gorm:"column:operator_id;" json:"operator_id"`
	RequestId          *big.Int  `gorm:"column:request_id;" json:"request_id"`
	WithdrawHeight     *big.Int  `gorm:"column:withdraw_height;" json:"withdraw_height"`
	WithdrawNethAmount string    `gorm:"column:withdraw_neth_amount;" json:"withdraw_neth_amount"`
	WithdrawExchange   string    `gorm:"column:withdraw_exchange;" json:"withdraw_exchange"`
	ClaimEthAmount     string    `gorm:"column:claim_eth_amount;" json:"claim_eth_amount"`
	Owner              string    `gorm:"column:owner;" json:"owner"`
	IsExit             bool      `gorm:"column:is_exit;" json:"is_exit"`
	Ctime              time.Time `gorm:"column:ctime;default" json:"ctime"`
	Mtime              time.Time `gorm:"column:mtime;default" json:"mtime"`
}

var nethWithdrawalRequests []*NethWithdrawalRequest

func NewNethWithdrawalRequest() *NethWithdrawalRequest {
	return new(NethWithdrawalRequest)
}

func (e *NethWithdrawalRequest) TableName() string {
	return "neth_withdrawal_request"
}

func (e *NethWithdrawalRequest) CreateBatch(withdrawalRequests []*NethWithdrawalRequest) error {
	return config.GlobalDB.Table(e.TableName()).Create(withdrawalRequests).Error
}

func (e *NethWithdrawalRequest) GetByRequestIds(network string, operatorId *big.Int, requestId []*big.Int, isExit bool) ([]*NethWithdrawalRequest, error) {
	db := config.GlobalDB.Table(e.TableName()).Where("network = ? AND operator_id = ? AND request_id IN ? AND type =? AND is_exit = ?", network, operatorId, requestId, isExit).Find(&nethWithdrawalRequests)
	return nethWithdrawalRequests, db.Error
}

func (e *NethWithdrawalRequest) UpdateExited(network string, requestId []*big.Int) error {
	return config.GlobalDB.Table(e.TableName()).Where("network = ? AND request_id IN ?", network, requestId).Updates(map[string]interface{}{"is_exit": true}).Error
}

func (e *NethWithdrawalRequest) InsertForNotExist(withdrawalRequests []*NethWithdrawalRequest) error {
	for _, request := range withdrawalRequests {
		db := config.GlobalDB.Table(e.TableName()).Where(NethWithdrawalRequest{RequestId: request.RequestId}).FirstOrCreate(&request)
		if db.Error != nil {
			return errors.Wrap(db.Error, "InsertForNotExist err.")
		}
	}
	return nil
}
