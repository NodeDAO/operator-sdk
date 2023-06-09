// description:
// @author renshiwei
// Date: 2023/6/5

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

func TestVnftExitScanByDB_Example_PreGenerateData_ByLocal(t *testing.T) {
	// 2023/6/9 need to exit tokenIds: [9 11 15 12 24 16]

	config.InitConfig("../../../conf/config-goerli.dev.yaml")
	//config.InitConfig("../../conf/config-mainnet.dev.yaml")
	logger.InitLog(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Format)
	err := config.InitOnce()
	require.NoError(t, err)

	nodedaoValidators := make([]*dao.NodedaoValidator, 0)

	nodedaoValidators = append(nodedaoValidators,
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0xae84b4a6d06bae763eccff9dfec868b6cfb180e092079c159b52bc7c70bc070b06fa242c03a1d4386bcd53f0e708fd25",
			OperatorId: 1,
			TokenId:    9,
			Type:       exitscan.VNFT,
		},
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0xa0092903b253492624ca0d60adb5edc41ff459c4b5a62b0b69565d88354e4ff1076845b038c968ea8ad92bf8c3ef2430",
			OperatorId: 1,
			TokenId:    11,
			Type:       exitscan.VNFT,
		},
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0xb488c66224ce62da3e429cc0a5baff3fb7814116c4f68eaf7b94f5278806b418f5f9b02b70e20a4e236a9268b5b4a5a0",
			OperatorId: 1,
			TokenId:    15,
			Type:       exitscan.VNFT,
		},
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0xae84b4a6d06bae763eccff9dfec868b6cfb180e092079c159b52bc7c70bc070b06fa242c03a1d4386bcd53f0e708fd25",
			OperatorId: 1,
			TokenId:    12,
			Type:       exitscan.VNFT,
		},
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0x834e9f350d42577220500ee4fbb238877b42f44293a6e1313ef5d8f7d4bd36f1c19a2e5c491c2629c9f46eae14dd09c7",
			OperatorId: 1,
			TokenId:    24,
			Type:       exitscan.VNFT,
		},
		&dao.NodedaoValidator{
			Network:    "goerli",
			Pubkey:     "0x8852f67c676b66abe8815f6555ba30aaf0f5a9775c6a59b96ea2693866c059c64b79de9b036dc62e0d192e47fdbea023",
			OperatorId: 1,
			TokenId:    16,
			Type:       exitscan.VNFT,
		},
	)

	nodedaoValidatorDao := dao.NewNodedaoValidator()
	err = nodedaoValidatorDao.CreateBatch(nodedaoValidators)
	require.NoError(t, err)
}

func TestVnftExitScanByDB_Example_ByLocal(t *testing.T) {
	ctx := context.Background()

	err := VnftExitScanByDB_Example(ctx)
	require.NoError(t, err)
}

func TestVnftExitScanByDB_Example_PostDeleteData_ByLocal(t *testing.T) {

}
