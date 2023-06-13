// description:
// @author renshiwei
// Date: 2023/6/13

package register

import (
	"context"
	"github.com/NodeDAO/operator-sdk/common/logger"
	"github.com/NodeDAO/operator-sdk/config"
	"github.com/NodeDAO/operator-sdk/contracts"
	"github.com/NodeDAO/operator-sdk/eth1"
	"github.com/NodeDAO/operator-sdk/example/dao"
	"github.com/NodeDAO/operator-sdk/validator/exitscan"
	"github.com/NodeDAO/operator-sdk/validator/register"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"math/big"
)

func RegisterValidatorForDB_Example(ctx context.Context) error {
	// init
	config.InitConfig("../../conf/config-goerli.dev.yaml")
	logger.InitLog(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Format)
	err := config.InitOnce()
	if err != nil {
		return errors.Wrap(err, "Failed to InitOnce.")
	}

	operatorId := big.NewInt(int64(config.GlobalConfig.Operator.Id))
	network := config.GlobalConfig.Eth.Network
	elAddr := config.GlobalConfig.Eth.ElAddr

	newRegisterValidator, err := register.NewRegisterValidator(ctx, network, elAddr, config.GlobalConfig.Operator.ControllerAddressPrivateKey, 40)
	if err != nil {
		return errors.Wrap(err, "Failed to RegisterValidator.")
	}

	//registerCount, err := newRegisterValidator.RegisterCount(operatorId.Int64())
	//if err != nil {
	//	return errors.Wrap(err, "Failed to RegisterCount.")
	//}
	registerCount := 1

	depositDatas := make([]*register.DepositData, 0, registerCount)
	depositDatas = append(depositDatas, &register.DepositData{
		Pubkey:                "",
		WithdrawalCredentials: "",
		Amount:                0,
		Signature:             "",
		DepositMessageRoot:    "",
		DepositDataRoot:       "",
		ForkVersion:           "",
		NetworkName:           "",
		DepositCliVersion:     "",
	})

	// register validator
	pubkeyBytes, err := newRegisterValidator.RegisterValidatorForNodeDAO(operatorId.Int64(), depositDatas)
	if err != nil {
		return errors.Wrap(err, "Failed to RegisterValidatorForNodeDAO.")
	}

	// -------- validator insert to db -----------
	elClient, err := eth1.NewEthClient(ctx, elAddr)
	if err != nil {
		return errors.Wrapf(err, "Fail to new eth client. network:%s", network)
	}
	vnftContract, err := contracts.NewVnftContract(network, elClient.Client)
	if err != nil {
		return errors.Wrapf(err, "Fail to new vnftContract. network:%s", network)
	}
	vnftOwnerVerify, err := exitscan.NewVnftOwnerVerify(ctx, network, elAddr)
	if err != nil {
		return errors.Wrap(err, "Failed to new VnftOwnerValidator.")
	}

	nodedaoValidators := make([]*dao.NodedaoValidator, 0, registerCount)
	for _, pubkeyByte := range pubkeyBytes {
		tokenId, err := vnftContract.Contract.TokenOfValidator(nil, pubkeyByte)
		if err != nil {
			return errors.Wrapf(err, "Fail to TokenOfValidator.")
		}
		stakeType, err := vnftOwnerVerify.VerifyStakeType(network, tokenId)
		if err != nil {
			return errors.Wrapf(err, "Fail to VerifyStakeType.")
		}
		nodedaoValidators = append(nodedaoValidators, &dao.NodedaoValidator{
			Network:    network,
			Pubkey:     hexutil.Encode(pubkeyByte),
			OperatorId: operatorId.Uint64(),
			TokenId:    tokenId.Uint64(),
			Type:       stakeType,
		})
	}

	nodedaoValidatorDao := dao.NewNodedaoValidator()
	err = nodedaoValidatorDao.CreateBatch(nodedaoValidators)

	return nil
}
