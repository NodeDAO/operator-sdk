// description:
// @author renshiwei
// Date: 2023/6/8

package dbscan

import (
	"context"
	"github.com/NodeDAO/operator-sdk/common/logger"
	"github.com/NodeDAO/operator-sdk/config"
	"github.com/NodeDAO/operator-sdk/validator/exitscan"
	"github.com/pkg/errors"
	"math/big"
)

// CornVnftExitScanByDB_Example corn vnft exit scan by db example
func CornVnftExitScanByDB_Example(ctx context.Context) {
	err := VnftExitScanByDB_Example(ctx)
	if err != nil {
		logger.Errorf("VnftExitScanByDB corn err:%+v", err)
	}
}

// VnftExitScanByDB_Example vnft exit scan by db example
func VnftExitScanByDB_Example(ctx context.Context) error {
	// init
	config.InitConfig("../../../conf/config-goerli.dev.yaml")
	logger.InitLog(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Format)
	err := config.InitOnce()
	if err != nil {
		return errors.Wrap(err, "Failed to InitOnce.")
	}

	// `config.GlobalConfig` is the global configuration information obtained by reading the configuration file
	// @see `operator-sdk/config`
	operatorId := big.NewInt(int64(config.GlobalConfig.Operator.Id))
	network := config.GlobalConfig.Eth.Network
	elAddr := config.GlobalConfig.Eth.ElAddr

	// 1. Scan the token ID of the v NFT that needs to be exited on the blockchain.
	vnftExitScan, err := exitscan.NewVnftExitScan(ctx, network, elAddr)
	if err != nil {
		return errors.Wrap(err, "Failed to NewVnftExitScan.")
	}

	contractVnftExitRecords, err := vnftExitScan.ExitScan(operatorId)
	if err != nil {
		return errors.Wrap(err, "Failed to ExitScan.")
	}

	// 2. Filter Validators that have initiated exits.
	vnftOwnerVerify, err := exitscan.NewVnftOwnerVerify(ctx, network, elAddr)
	if err != nil {
		return errors.Wrap(err, "Failed to new VnftOwnerValidator.")
	}
	dbFilter, err := NewDBFiltrate(network, exitscan.VNFT, vnftOwnerVerify)
	if err != nil {
		return errors.Wrap(err, "Failed to new DBFilter.")
	}
	filterVnftExitRecords, err := dbFilter.Filter(operatorId, contractVnftExitRecords)
	if err != nil {
		return errors.Wrap(err, "Failed to Filter by DB.")
	}

	// 3. Exit Validator and mark it in the database.
	// !!!!!!! Only the markup is made here, there is no real exit from the Validator,
	// !!!!!!! and the operation of exiting the Validator needs to be implemented
	// !!!!!!! by the operator according to its own technology stack.
	dbMark, err := NewDBMark(network)
	if err != nil {
		return errors.Wrap(err, "Failed to new NewDBMark.")
	}
	err = dbMark.ExitMark(operatorId, filterVnftExitRecords)
	if err != nil {
		return errors.Wrap(err, "Failed to ExitMark by db.")
	}

	return nil
}
