// description:
// @author renshiwei
// Date: 2023/6/5

package exitscan

import (
	"context"
	"fmt"
	"github.com/NodeDAO/operator-sdk/config"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestNethExitScan_ByLocal(t *testing.T) {
	ctx := context.Background()
	config.InitConfig("../../conf/config-goerli.dev.yaml")
	//config.InitConfig("../../conf/config-mainnet.dev.yaml")

	nethExitScan, err := NewNETHExitScan(ctx, config.GlobalConfig.Eth.Network, config.GlobalConfig.Eth.ElAddr)
	require.NoError(t, err)

	withdrawalInfo, err := nethExitScan.WithdrawalRequestScan(big.NewInt(1))
	require.NoError(t, err)
	fmt.Printf("withdrawalInfo scan:%+v\n", withdrawalInfo)

	nethRecords, err := nethExitScan.ExitScan(big.NewInt(1))
	require.NoError(t, err)

	fmt.Printf("neth exit scan:%+v\n", nethRecords)
}
