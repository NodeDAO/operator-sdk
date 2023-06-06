// description:
// @author renshiwei
// Date: 2023/6/6

package dao

import (
	"math/big"
	"time"
)

type NethWithdrawalRequest struct {
	ID                 uint64    `gorm:"column:id;primary_key" json:"id"`
	Network            string    `gorm:"column:network;" json:"network"`
	OperatorId         *big.Int  `gorm:"column:operator_id;" json:"operator_id"`
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
