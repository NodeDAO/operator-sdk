# vNFT exitscan steps
1. **Scan the tokenId of the vNFT that needs to be exited on the chain**.
This can be done through the withdrawalRequest contract `getUserUnstakeButOperatorNoExitNfs` method.

This process filters the tokenId, and further filtering is required, because the state of the Validator that exits asynchronously may not yet be synchronized to the chain.
2. **Operator implements Validator filtering that has initiated exit**.

The exit of the Validator needs to go through the lifetime of the beacon, the exit is asynchronous, and the exit completion time is not certain. The Validator that initiates the exit does not report to Oracle to trigger settlement, and the on-chain data will not be updated. Therefore, off-chain means are needed to filter.

Different operators may have different technical implementations, and we provide the simplest example for reference, using MySQL for filtering.

3. **Exit Validator and mark it in the database**.

> The above steps need to be performed through a scheduled task cycle.

# Specific implementation

## 1. Scan the vNFT on the chain to be exited

For implementation, see: [validator/exitscan/vnft.go](../../validator/exitscan/vnft.go)

Method: `ExitScan`

## 2. Filter the Validators that have initiated the exit

> This part of the example Here with the help of MySQL implementation, the operator can be implemented according to the specific situation.
>
> Before filtering exiting Validator, the Validator should already exist in the database table `nodedao_validator`. You can then `registerValidator` and store it in the library, which is explained in the `registerValidator` operation. It can now be assumed that the Validator is already in table `nodedao_validator`.

For example: [example/exitscan/dbscan/exit_filter.go](../../example/exitscan/dbscan/exit_filter.go)

Method: `Filter`

`VnftExitScan.ExitScan` scans the on-chain record `vnftRecords` as the input parameter `vnftContractExitRecords` of `Filter`.

## 3. Exit the Validator and mark it in the database

> This part of the example Here with the help of MySQL implementation, the operator can be implemented according to the specific situation.

In the second step, the filtered validators are what the operator needs to exit, and the operator needs to initiate the exit of the validator according to its own technical implementation. Mark these validators as having exited for the next `filter`.

For example, see: [example/exitscan/dbscan/exit_mark.go](../../example/exitscan/dbscan/exit_mark.go)

Method: `ExitMark`


## Example implementation process

For code, see: [example/exitscan/dbscan/vnft_exit_scan_example.go](../../example/exitscan/dbscan/vnft_exit_scan_example.go)

Method: `VnftExitScanByDB_Example`

> In the code of step 3, you need to exit the filtered validator, and this part of the operator needs to be implemented according to its own technology stack.

## Scheduled tasks are executed periodically

The above three steps should be performed periodically in a timed task, a user initiates unstake, and the operator should exit the Validator in time.

For methods called by scheduled tasks, see: Code See: [example/exitscan/dbscan/vnft_exit_scan_example.go](../../example/exitscan/dbscan/vnft_exit_scan_example.go)

Method: `CornVnftExitScanByDB_Example`

> Timed task implementation: omitted