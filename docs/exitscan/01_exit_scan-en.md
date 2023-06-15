# NodeDAO Unstake
The NodeDAO protocol supports users to unstake:

1. vNFT directly initiates unstake
2. Greater than 32nETH supports largeRequest, and unstake is initiated through largeRequest; If it is less than 32nETH, you need to use Operator's LiquidityPool for unstake.

**For vNFT unstake and nETH largeRequest, the operator needs to actively exit the corresponding validator**.



# Why do you need exit scan?
When the user initiates the unstake or largeRequest, the corresponding data will be recorded on the chain, but the data on the chain will not actively notify the Operator, and the Operator needs to actively query the data on the chain and then perform the corresponding operation. This process requires exit scan.

> This section can be found at: [validator/exitscan](../../validator/exitscan)

Because the exit of the Validator is asynchronous and needs to go through the lifetime of the beacon, during the time of the beacon exit cycle, the on-chain smart contract does not know that these Validators have initiated the exit. So there are other means to mark these validators.

> the most intuitive way to use databases, we implemented an example with MySQL: [example/exitscan/dbscan](../../example/exitscan/dbscan)
> Operator has different technical implementations, this section is for reference only.

It is necessary to wait for Oracle to report the Validator exit and trigger the settlement before these exited validator information will be updated to the chain.



# exit scan abstraction layer

See Abstraction layer code：[validator/exitscan/typings.go](../../validator/exitscan/typings.go)

## VnftOwner and StakeType

```go
type VnftOwner uint32
type StakeType uint32

const (
	USER VnftOwner = iota
	LiquidStaking
)

const (
	VNFT StakeType = iota
	NETH
)

func GetVnftOwner(stakeType StakeType) VnftOwner {
	if stakeType == VNFT {
		return USER
	} else {
		return LiquidStaking
	}
}
```

**Each validator corresponds to a vNFT, represented by a unique ID tokenId. Go to stake with vNFT to get vNFT held by the user; to obtain nETH staking, the vNFT corresponding to the 32ETH validator behind it is owned by LiquidStaking.** 


## VnftRecord

code：

```go
type VnftRecord struct {
	Network    string
	OperatorId *big.Int
	TokenId    *big.Int
	Pubkey     string
	Type       StakeType
}
```

`VnftRecord` Used to indicate that a Validator record needs to be exited。



## VnftOwnerValidator

code：

```go
type VnftOwnerValidator interface {
	// VerifyVnftOwner Verify that the stakeType of vnft tokenIds and vnftOwner match
	// ----------------------------------------------------------------
	// The relationship between StakeType and VnftOwner is as follows:
	// ----------------------
	// StakeType | VnftOwner
	// ----------------------
	// VNFT      | USER
	// NETH      | LiquidStaking
	VerifyVnftOwner(network string, stakeType StakeType, vnftOwner VnftOwner, tokenIds []*big.Int) (bool, error)
}
```

`VnftOwnerValidator` is used to verify that the relationship between the validator `StakeType` and `VnftOwner` is correct.


## ExitScanner

code：

```go
// ExitScanner Scan the smart contract for records that need to be exited
type ExitScanner interface {
	ExitScan(operatorId *big.Int) ([]*VnftRecord, error)
}
```

`ExitScanner` The function is mainly to scan the Validator from the chain that needs to be exited。

Operators can use their methods directly.

vNFT needs to implement `ExitScanner` and nETHLargeRequest needs to implement `WithdrawalRequestScanner`. You also need to use the `ExitFilter` to filter, and finally get the final list of validators that need to be exited.


## WithdrawalRequestScanner

code：

```go
// WithdrawalRequestScanner nETH exit depends on the WithdrawalRequest.
// vNFT exit can be used directly with Exit Scanner
type WithdrawalRequestScanner interface {
	ExitScanner

	// WithdrawalRequestScan Scan for unclaimed Withdrawal Requests
	WithdrawalRequestScan(operatorId *big.Int) ([]*withdrawalRequest.WithdrawalRequestWithdrawalInfo, error)
}
```

`WithdrawalRequestScanner` inherits from `ExitScanner` and extends the `WithdrawalRequestScan` method.
Operators can use their methods directly.

nETHLargeRequest scans the list of unprocessed largeRequests on the chain via the `WithdrawalRequestScan` method.

Scan the chain for validators that need to be exited via `ExitScanner`. You also need to use the `ExitFilter` to filter, and finally get the final list of validators that need to be exited.


## ExitFilter

The exit of the Validator needs to go through the lifetime of the beacon, the exit is asynchronous, and the exit completion time is not certain. The Validator that initiates the exit does not report to Oracle to trigger settlement, and the on-chain data will not be updated. Therefore, off-chain means are needed to filter.
code：

```go
// ExitFilter filter the exit for vNFT and nETH.
// Validator's exit is asynchrony. The reasons for asynchrony are:
// 1. The validator exit goes through the lifetime of the beacon.
// 2. NodeDAO-Oracle is required to report to settle.
// --------------------------------------------
// Filter The operator needs to implement it by itself, and the easiest way is to use db.
// An example implementation will be provided, based on MySQL, see Example.
type ExitFilter interface {
	// Filter To filter for exit
	// @return []*VnftRecord{} Filtered
	Filter(operatorId *big.Int, vnftRecords []*VnftRecord) ([]*VnftRecord, error)
}
```

**This part of the operator can be implemented in combination with its own technology to implement the `ExitFilter` interface.** 

`Filter` will pass in a `[]*VnftRecord` and return a `[]*VnftRecord`, and the intermediate operation filters the exited validator.

> different operators may have different technical implementations, we provide one of the simplest examples for reference, using MySQL to implement filtering. Example: [example/exitscan/dbscan](../../example/exitscan/dbscan)
>
> **This part of the operator can be implemented in combination with its own technology to implement the 'ExitFilter' interface.** 


## WithdrawalRequestFilter

code：

```go
// WithdrawalRequestFilter To filter for WithdrawalRequests.
// --------------------------------------------
// The simplest way to implement the operator is to use db, see example.
type WithdrawalRequestFilter interface {
	// WithdrawalRequestFilter To filter for WithdrawalRequests.
	// @return []*WithdrawalRequest Filtered WithdrawalRequests.
	WithdrawalRequestFilter(operatorId *big.Int, withdrawalRequests []*WithdrawalRequest) ([]*WithdrawalRequest, error)
}
```

`WithdrawalRequestFilter` is used to handle nETHLargeRequest, as it is also asynchronous and needs to be processed off-chain.
> different operators may have different technical implementations, we provide one of the simplest examples for reference, using MySQL to implement filtering. Example: [example/exitscan/dbscan](../../example/exitscan/dbscan)
>
> **This part of the operator can be implemented in combination with its own technology to implement the `WithdrawalRequestFilter` interface.** 



## WithdrawalRequestExitValidatorCounter

code：

```go
// WithdrawalRequestExitValidatorCounter Calculate the number of validators that need to be exited by a Withdrawal Request
type WithdrawalRequestExitValidatorCounter interface {
	// ExitCounter Calculate the number of validators that need to be exited by a Withdrawal Request
	// @param filterWithdrawalRequests  A list of offline filtered Withdrawal Requests
	ExitCounter(filterWithdrawalRequests []*WithdrawalRequest) (uint32, error)
}
```

The implementation code can be seen in ：[validator/exitscan/neth_large_request.go](../../validator/exitscan/neth_large_request.go)

The number of exit validators for nETH largeRequests is calculated primarily from the filtered `filterWithdrawalRequests` list.
The brief calculation logic is to calculate the number of validators exiting based on calculating the sum of `ClaimEthAmount` for `filterWithdrawalRequests`.
> Evaluate the rule：
>
> if sumETHAmount = 64 ether, need to exit 2 validator
> if sumETHAmount = 66 ether, need to exit 3 validator



## ExitMarker

code：

```go
// ExitMarker To perform a validator exit, it needs to be flagged, and then it is used for filter
// --------------------------------------------
// The simplest way to implement the operator is to use db, see example
type ExitMarker interface {
	// ExitMark Mark the exit of the Vnft Record
	ExitMark(operatorId *big.Int, records []interface{}) error
}
```

When the Validator initiates an exit, mark it as exited so that the Filter can filter it.

>  Different operators may have different technical implementations, and we provide one of the simplest examples for reference, using My SQL to implement filtering. example：[example/exitscan/dbscan](../../example/exitscan/dbscan)
>
> **This part of the operator can be implemented in combination with its own technology to implement the 'ExitMarker' interface.**



## WithdrawalRequestMarker

The `WithdrawalRequest` also needs to be flagged after processing for the next filter.

code：

```go
// WithdrawalRequestMarker To mark deal for WithdrawalRequest.
// --------------------------------------------
// The simplest way to implement the operator is to use db, see example.
type WithdrawalRequestMarker interface {
	// WithdrawalRequestMark To mark deal for WithdrawalRequest
	WithdrawalRequestMark(operatorId *big.Int, withdrawalRequests []*WithdrawalRequest) error
}
```

> different operators may have different technical implementations, we provide one of the simplest examples for reference, using MySQL to implement filtering. Example: [example/exitscan/dbscan](../../example/exitscan/dbscan)
>
> **This part of the operator can be implemented in combination with its own technology to implement the `WithdrawalRequestMarker` interface.** 





# Example by MySQL
The examples of `ExitFilter` and `ExitMarker` use MySQL.

SQL file see: [script/sql/operator-sdk-*.sql](../../script/sql/operator-sdk-*.sql)

- `nodedao_validator` to store the necessary information of the Validator (including: tokenId, operatorId, whether to initiate an exit, owner, etc.)
- `neth_withdrawal_request` is used to store nETH largeRequest information and mark whether to process exit information.



