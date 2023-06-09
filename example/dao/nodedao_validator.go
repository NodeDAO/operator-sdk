// description:
// @author renshiwei
// Date: 2023/6/5

package dao

import (
	"github.com/NodeDAO/operator-sdk/config"
	"github.com/NodeDAO/operator-sdk/validator/exitscan"
	"time"
)

type NodedaoValidator struct {
	ID         uint64             `gorm:"column:id;primary_key" json:"id"`
	Network    string             `gorm:"column:network;" json:"network"`
	Pubkey     string             `gorm:"column:pubkey;" json:"pubkey"`
	OperatorId uint64             `gorm:"column:operator_id;type:BIGINT;" json:"operator_id"`
	TokenId    uint64             `gorm:"column:token_id;type:BIGINT;" json:"token_id"`
	Type       exitscan.StakeType `gorm:"column:type;" json:"type"`
	IsExit     bool               `gorm:"column:is_exit;" json:"is_exit"`
	Ctime      time.Time          `gorm:"column:ctime;default" json:"ctime"`
	Mtime      time.Time          `gorm:"column:mtime;default" json:"mtime"`
}

var nodedaoValidators []*NodedaoValidator

func NewNodedaoValidator() *NodedaoValidator {
	return new(NodedaoValidator)
}

func (e *NodedaoValidator) TableName() string {
	return "nodedao_validator"
}

func (e *NodedaoValidator) CreateBatch(vnftRecords []*NodedaoValidator) error {
	return config.GlobalDB.Table(e.TableName()).Create(vnftRecords).Error
}

func (e *NodedaoValidator) GetByTokenIds(network string, operatorId uint64, tokenIds []uint64, unstakeType exitscan.StakeType, isExit bool) ([]*NodedaoValidator, error) {
	db := config.GlobalDB.Table(e.TableName()).Where("network = ? AND operator_id = ? AND token_id IN ? AND type =? AND is_exit = ?", network, operatorId, tokenIds, unstakeType, isExit).Find(&nodedaoValidators)
	return nodedaoValidators, db.Error
}

func (e *NodedaoValidator) UpdateExited(network string, tokenIds []uint64) error {
	return config.GlobalDB.Table(e.TableName()).Where("network = ? AND token_id IN ?", network, tokenIds).Updates(map[string]interface{}{"is_exit": true}).Error
}
