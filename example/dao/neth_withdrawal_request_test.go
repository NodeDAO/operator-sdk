// description:
// @author renshiwei
// Date: 2023/6/16

package dao

import (
	"fmt"
	"github.com/NodeDAO/operator-sdk/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNethWithdrawalRequest_InsertForNotExist(t *testing.T) {
	config.InitConfig("../../conf/config-goerli.dev.yaml")
	err := config.InitOnce()
	require.NoError(t, err)

	nethWithdrawalRequestDao := NewNethWithdrawalRequest()
	withdrawalRequests := make([]*NethWithdrawalRequest, 0)
	err = nethWithdrawalRequestDao.InsertForNotExist(withdrawalRequests)
	require.NoError(t, err)

	requestIds := make([]uint64, 0)
	requestIds = append(requestIds, 2)
	withdrawalRequestInfos, err := nethWithdrawalRequestDao.GetByRequestIds("goerli", 1, requestIds, false)
	require.NoError(t, err)

	fmt.Println(withdrawalRequestInfos)
}
