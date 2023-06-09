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

// CornNethExitScanByDB_Example corn nETH exit scan by db example
func CornNethExitScanByDB_Example(ctx context.Context) {
	err := NethExitScanByDB_Example(ctx)
	if err != nil {
		logger.Errorf("VnftExitScanByDB corn err:%+v", err)
	}
}

// NethExitScanByDB_Example nETH large request exit scan by db example.
func NethExitScanByDB_Example(ctx context.Context) error {
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

	exitScan, err := exitscan.NewNETHExitScan(ctx, network, elAddr)
	if err != nil {
		return errors.Wrap(err, "Failed to NewVnftExitScan.")
	}

	// 1. Scan the list of nETH largeRequestId on the chain that needs to be exited.
	withdrawalInfos, err := exitScan.WithdrawalRequestScan(operatorId)
	if err != nil {
		return errors.Wrap(err, "Failed to WithdrawalRequestScan.")
	}

	// 2. Get the record of Operator nETH corresponding to vNFT on the chain.
	contractVnftExitRecordsOfNeth, err := exitScan.ExitScan(operatorId)
	if err != nil {
		return errors.Wrap(err, "Failed to ExitScan.")
	}

	// 3. Filter the list of withdrawal requests.
	vnftOwnerVerify, err := exitscan.NewVnftOwnerVerify(ctx, network, elAddr)
	if err != nil {
		return errors.Wrap(err, "Failed to new VnftOwnerValidator.")
	}
	dbFilter, err := NewDBFiltrate(network, exitscan.NETH, vnftOwnerVerify)
	if err != nil {
		return errors.Wrap(err, "Failed to new DBFilter.")
	}
	dbFilter.SetExitValidatorCounter(exitScan)

	filterWithdrawalRequests, err := dbFilter.WithdrawalRequestFilter(operatorId, withdrawalInfos)
	if err != nil {
		return errors.Wrap(err, "Failed to Filter by DB.")
	}

	// 4. Filter Validators that have initiated exits.
	filterVnftExitRecords, err := dbFilter.Filter(operatorId, contractVnftExitRecordsOfNeth)
	if err != nil {
		return errors.Wrap(err, "Failed to Filter by DB.")
	}

	// 5. Calculate how many validators exit based on the filtered withdrawal request list.
	// This part of the implementation has been completed in steps 3 and 4.

	// 6. Exit Validator and mark it in the database.
	// !!!!!!! Only the markup is made here, there is no real exit from the Validator,
	// !!!!!!! and the operation of exiting the Validator needs to be implemented
	// !!!!!!! by the operator according to its own technology stack.
	dbMark, err := NewDBMark(network)
	if err != nil {
		return errors.Wrap(err, "Failed to new NewDBMark.")
	}
	err = dbMark.ExitMark(operatorId, filterVnftExitRecords)
	if err != nil {
		config.GlobalDB.Rollback()
		return errors.Wrap(err, "Failed to ExitMark by db.")
	}

	// 7. Mark the filtered withdrawal request list as processed.
	err = dbMark.WithdrawalRequestMark(operatorId, filterWithdrawalRequests)
	if err != nil {
		return errors.Wrap(err, "Failed to WithdrawalRequestMark by db.")
	}

	return nil
}
