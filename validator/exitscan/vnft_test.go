// description:
// @author renshiwei
// Date: 2023/6/5

package exitscan

import (
	"context"
	"fmt"
	"github.com/NodeDAO/operator-sdk/common/logger"
	"github.com/NodeDAO/operator-sdk/config"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestVnftExitScan_ByLocal(t *testing.T) {
	ctx := context.Background()
	config.InitConfig("../../conf/config-goerli.dev.yaml")
	//config.InitConfig("../../conf/config-mainnet.dev.yaml")
	logger.InitLog(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Format)

	vnftExitScan, err := NewVnftExitScan(ctx, config.GlobalConfig.Eth.Network, config.GlobalConfig.Eth.ElAddr)
	require.NoError(t, err)

	vnftRecords, err := vnftExitScan.ExitScan(big.NewInt(1))
	require.NoError(t, err)

	fmt.Printf("vnft exit scan:%+v\n", vnftRecords)
}
