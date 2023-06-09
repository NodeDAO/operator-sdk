// description:
// @author renshiwei
// Date: 2023/6/6

package dbscan

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNethExitScanByDB_Example_PreGenerateData_ByLocal(t *testing.T) {
	// 2023/6/9
}

func TestNethExitScanByDB_Example_ByLocal(t *testing.T) {
	ctx := context.Background()

	err := NethExitScanByDB_Example(ctx)
	require.NoError(t, err)
}

func TestNethExitScanByDB_Example_PostDeleteData_ByLocal(t *testing.T) {

}
