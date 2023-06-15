# nETH large Request exitscan step

1. **List of nETH largeRequestId that need to be exited on the scan chain**. 

This can be done through the withdrawalRequest contract 'getWithdrawalOfRequestId' method.

This process filters the largeRequestId and needs further filtering, because the state of the Validator that exits asynchronously may not yet be synchronized to the chain.

2. **Get the record of Operator nETH corresponding to vNFT on the chain**.

**nETH corresponds to vNFT owned by LiquidStaking, nETH largeRequest to exit validator must be owned by LiquidStaking. **

> vNFT way records, the owner is the user.

3. **The operator implements the filter withdrawal request list by itself**. 

4. **Operator implements Validator filtering that has initiated exit**.

The exit of the Validator needs to go through the lifetime of the beacon, the exit is asynchronous, and the exit completion time is not certain. The Validator that initiates the exit does not report to Oracle to trigger settlement, and the on-chain data will not be updated. Therefore, off-chain means are needed to filter.

Different operators may have different technical implementations, and we provide the simplest example for reference, using MySQL for filtering.

5. **Calculate how many validators are exited according to the filtered withdrawalRequest list**.
6. **Exit Validator and mark it in the database**.
7. **Mark the filtered withdrawalRequest list as processed**.



# Specific implementation

## 1. Scan the list of nETH largeRequestId on the chain that needs to be exited

For implementation, see: [validator/exitscan/neth_large_request.go](../../validator/exitscan/neth_large_request.go)

Method: `WithdrawalRequestScan`

This method scans the chain for WithdrawalRequests and de-chain filtering.




## 2. Obtain the record of Operator nETH corresponding to vNFT on the chain

For implementation, see: [validator/exitscan/neth_large_request.go](../../validator/exitscan/neth_large_request.go)

Method: `ExitScan`

This method scans vNFT records that are not exited on-chain and owned by LiquidStaking, and also requires de-off-chain filtering.



## 3. the operator implements the filter withdrawal request list by itself

> This part of the example Here with the help of MySQL implementation, the operator can be implemented according to the specific situation.

For example: [example/exitscan/dbscan/exit_filter.go](../../example/exitscan/dbscan/exit_filter.go)

Method: `WithdrawalRequestFilter`

In the example, the number of validators that need to exit is calculated based on the filtered list of `filterWithdrawalRequests`.

For code, see the implementation of `WithdrawalRequestExitValidatorCounter`:

```go
// WithdrawalRequestExitValidatorCounter Calculate the number of validators that need to be exited by a Withdrawal Request
type WithdrawalRequestExitValidatorCounter interface {
	// ExitCounter Calculate the number of validators that need to be exited by a Withdrawal Request
	// @param filterWithdrawalRequests  A list of offline filtered Withdrawal Requests
	ExitCounter(filterWithdrawalRequests []*WithdrawalRequest) (uint32, error)
}
```



In addition, to implement `WithdrawalRequestFilter` off-chain filtering, there must be a way to store scanned `WithdrawalRequest` records off-chain. This part of the example is also implemented in the `WithdrawalRequestFilter`.

> Validator(vNFT) records store the records in the sample's database in the registerValidator section.



## 4. Operator implements Validator filtering that has initiated exit

> This part of the example Here with the help of MySQL implementation, the operator can be implemented according to the specific situation.
>
> Before filtering exiting Validator, the Validator should already exist in the database table `nodedao_validator`. You can then `registerValidator` and store it in the library, which is explained in the `registerValidator` operation. It can now be assumed that the Validator is already in table `nodedao_validator`.

For example: [example/exitscan/dbscan/exit_filter.go](../../example/exitscan/dbscan/exit_filter.go)

Method: `WithdrawalRequestFilter`

Method: `Filter`

`NETHExitScan.ExitScan` scans the on-chain record `vnftRecords` as the input parameter `vnftContractExitRecords` for `Filter`.

## 5. Calculate how many validators exit according to the filtered withdrawalRequest list

For implementation, see: [validator/exitscan/neth_large_request.go](../../validator/exitscan/neth_large_request.go)

Method: `ExitCounter`

> if sumETHAmount = 64 ether, need to exit 2 validator
> if sumETHAmount = 66 ether, need to exit 3 validator




## 6, exit Validator and mark it in the database

> This part of the example Here with the help of MySQL implementation, the operator can be implemented according to the specific situation.

In the second step, the filtered validators are what the operator needs to exit, and the operator needs to initiate the exit of the validator according to its own technical implementation. Mark these validators as having exited for the next `filter`.

For example, see: [example/exitscan/dbscan/exit_mark.go](../../example/exitscan/dbscan/exit_mark.go)

Method: `ExitMark`


## 7. Mark the filtered withdrawalRequest list as processed

> This part of the example Here with the help of MySQL implementation, the operator can be implemented according to the specific situation.

Mark these withdrawalRequests that have been processed for the next `Filter`.

See for an example：[example/exitscan/dbscan/exit_mark.go](../../example/exitscan/dbscan/exit_mark.go)

Method：`WithdrawalRequestMark`



## Example implementation process

See the code：[example/exitscan/dbscan/neth_exit_scan_example.go](../../example/exitscan/dbscan/neth_exit_scan_example.go)

Method：``NethExitScanByDB_Example`
> In the code of step 6, you need to exit the filtered validator, which needs to be implemented according to its own technology stack.



## Scheduled tasks are executed periodically

The above three steps should be performed periodically in a timed task, a user initiates unstake, and the operator should exit the Validator in time.

For methods called by scheduled tasks, see: Code See: [example/exitscan/dbscan/vnft_exit_scan_example.go](../../example/exitscan/dbscan/vnft_exit_scan_example.go)

Method: `CornNethExitScanByDB_Example`

> Timed task implementation: omitted

# User claim withdrawRequest

When the operator exits the validator for withdrawRequest, and the Oracle reports the triggering settlement, ETH returns to the Operator Liquidity Pool.

Users also need to make a claim for ETH to be credited to the withdrawal account. The claim prerequisite is: `Withdraw Amount <= Operator Liquidity Pool`.


![neth claim](../images/neth-claim.jpeg)

