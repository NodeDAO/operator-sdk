// description:
// @author renshiwei
// Date: 2023/6/6

package dbscan

import (
	"context"
	"github.com/NodeDAO/operator-sdk/common/logger"
	"github.com/NodeDAO/operator-sdk/config"
	"github.com/NodeDAO/operator-sdk/example/dao"
	"github.com/NodeDAO/operator-sdk/validator/exitscan"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNethExitScanByDB_Example_PreGenerateData_ByLocal(t *testing.T) {
	// 2023/6/9 requestId [0,2,7] exitValidatorCount: 4
	config.InitConfig("../../../conf/config-goerli.dev.yaml")
	//config.InitConfig("../../conf/config-mainnet.dev.yaml")
	logger.InitLog(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Format)
	err := config.InitOnce()
	require.NoError(t, err)

	// insert withdraw request
	//nethWithdrawalRequests := make([]*dao.NethWithdrawalRequest, 0)
	//nethWithdrawalRequests = append(nethWithdrawalRequests,
	//	&dao.NethWithdrawalRequest{
	//		Network:            "goerli",
	//		OperatorId:         1,
	//		RequestId:          0,
	//		WithdrawHeight:     "8851198",
	//		WithdrawNethAmount: "33000000000000000000",
	//		WithdrawExchange:   "1049483201626165037",
	//		ClaimEthAmount:     "34632945653663446248",
	//		Owner:              "0xF5ade6B61BA60B8B82566Af0dfca982169a470Dc",
	//	},
	//	&dao.NethWithdrawalRequest{
	//		Network:            "goerli",
	//		OperatorId:         1,
	//		RequestId:          2,
	//		WithdrawHeight:     "8898343",
	//		WithdrawNethAmount: "33000000000000000000",
	//		WithdrawExchange:   "1050977570362921211",
	//		ClaimEthAmount:     "34682259821976399981",
	//		Owner:              "0x3535d10Fc0E85fDBC810bF828F02C9BcB7C2EBA8",
	//	},
	//	&dao.NethWithdrawalRequest{
	//		Network:            "goerli",
	//		OperatorId:         1,
	//		RequestId:          7,
	//		WithdrawHeight:     "8980169",
	//		WithdrawNethAmount: "33481080000000000000",
	//		WithdrawExchange:   "1067485358872657139",
	//		ClaimEthAmount:     "35740562699244143484",
	//		Owner:              "0x00dFaaE92ed72A05bC61262aA164f38B5626e106",
	//	},
	//)
	//withdrawalRequestDao := dao.NewNethWithdrawalRequest()
	//err = withdrawalRequestDao.InsertForNotExist(nethWithdrawalRequests)
	//require.NoError(t, err)

	// insert nodedao validator
	nodedaoValidators := make([]*dao.NodedaoValidator, 0)
	nodedaoValidators = append(nodedaoValidators,
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0xad4aee13a02602472d40af69e723c004d6a057e9f2436b0eb55a7608ef1e4905a585f4a1e1a25a62da17f5da1921af71",
			OperatorId: 1,
			TokenId:    21,
			Type:       exitscan.NETH,
		},
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0x89cdbff38e2fee90095674c7eb91aebb5302a101811c588b31ee6c2f6555b44b43986f9b7a037a9c6fc94593e701106a",
			OperatorId: 1,
			TokenId:    22,
			Type:       exitscan.NETH,
		},
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0xae9d10a4ee2316ba0332e3f60bb2203d496b407de3fa5b7961e72dece63d7900cda67e61722bac061c6907c6035b30d9",
			OperatorId: 1,
			TokenId:    25,
			Type:       exitscan.NETH,
		},
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0xa5a24e7a142664d40949f3c4a9eb4d3aa028e74b4788ab6c385d512ea622d070db932d990eb7118f45a294d8aec026ed",
			OperatorId: 1,
			TokenId:    26,
			Type:       exitscan.NETH,
		},
	)

	nodedaoValidatorDao := dao.NewNodedaoValidator()
	err = nodedaoValidatorDao.CreateBatch(nodedaoValidators)
	require.NoError(t, err)

}

func TestNethExitScanByDB_Example_ByLocal(t *testing.T) {
	ctx := context.Background()

	err := NethExitScanByDB_Example(ctx)
	require.NoError(t, err)
}

func TestNethExitScanByDB_Example_PostDeleteData_ByLocal(t *testing.T) {

}
