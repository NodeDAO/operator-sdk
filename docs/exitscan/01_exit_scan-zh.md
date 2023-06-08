# NodeDAO Unstake
NodeDAO协议支持用户进行unstake：
1. vNFT直接发起unstake
2. 32nETH支持largeRequest，通过largeRequest发起unstake；如果小于32nETH，则需要使用Operator的LiquidityPool进行unstake。

**对于vNFT的unstake和nETH的largeRequest，Operator需要主动去退出对应的validator**。



# 为什么需要 exit scan？
当用户发起发起unstake或largeRequest后，相应的数据会被记录在链上，但是链上的数据并不会主动通知Operator，Operator需要主动去查询链上的数据，然后进行相应的操作。这个过程需要exit scan。

> 这部分内容可以参看：[validator/exitscan](../../validator/exitscan)



因为Validator的退出是异步的，需要经过beacon的生命周期，在beacon退出周期的时间内，链上智能合约并不知道这些Validator已发起退出。所以还需要借助其他手段对这些Validator进行标记。

> 最直观的方式使用数据库， 我们通过MySQL实现了一个示例：[example/exitscan/dbscan](../../example/exitscan/dbscan)
> Operator拥有不同的技术实现，这部分内容仅供参考。



需要等待Oracle上报Validator退出并触发结算后，这些退出的validator信息才会更新到链上。



# exit scan 抽象层

抽象层代码参看：[validator/exitscan/typings.go](../../validator/exitscan/typings.go)

## VnftOwner和StakeType

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

**每一个Validator都对应一个vNFT，由唯一ID tokenId所表示。以vNFT去stake获得vNFT由用户所持有；以获得nETH的质押，背后代表32ETH validator对应的vNFT由LiquidStaking所有。**



## VnftRecord

代码：

```go
type VnftRecord struct {
	Network    string
	OperatorId *big.Int
	TokenId    *big.Int
	Pubkey     string
	Type       StakeType
}
```

`VnftRecord` 用于表示需要退出的Validator记录。



## VnftOwnerValidator

代码：

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

`VnftOwnerValidator` 用来验证validator `StakeType` 和 `VnftOwner` 的关系是否正确。



## ExitScanner

代码：

```go
// ExitScanner Scan the smart contract for records that need to be exited
type ExitScanner interface {
	ExitScan(operatorId *big.Int) ([]*VnftRecord, error)
}
```

`ExitScanner` 的功能主要是从链上扫描出需要退出的Validator。

Operator可以直接使用其方法。

vNFT需要实现`ExitScanner`，nETHLargeRequest需要实现 `WithdrawalRequestScanner`。还需要使用 `ExitFilter` 去过滤，最后获得最终的需要退出的Validator列表。 



## WithdrawalRequestScanner

代码：

```go
// WithdrawalRequestScanner nETH exit depends on the WithdrawalRequest.
// vNFT exit can be used directly with Exit Scanner
type WithdrawalRequestScanner interface {
	ExitScanner

	// WithdrawalRequestScan Scan for unclaimed Withdrawal Requests
	WithdrawalRequestScan(operatorId *big.Int) ([]*withdrawalRequest.WithdrawalRequestWithdrawalInfo, error)
}
```

`WithdrawalRequestScanner` 继承与 `ExitScanner`，并扩展了 `WithdrawalRequestScan` 方法。

Operator可以直接使用其方法。

nETHLargeRequest通过 `WithdrawalRequestScan` 方法去扫描链上未处理的 largeRequest列表。

通过 `ExitScanner` 扫描链上需要退出的Validator。还需要使用 `ExitFilter` 去过滤，最后获得最终的需要退出的Validator列表。 



## ExitFilter

Validator的退出需要经过beacon的生命周期，退出是异步的，退出完成时间并不确定。发起退出的Validator在Oracle没有上报触发结算，链上数据是并不会更新。因此需要链下的手段来进行过滤。

代码：

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

**这部分Operator可结合自己的技术进行实现，实现 `ExitFilter` 接口。**

`Filter` 会传入一个 `[]*VnftRecord`，返回一个 `[]*VnftRecord`，中间操作过滤已发起退出的validator。

> 不同Operator可能存在不同的技术实现，我们提供一种最简单示例用以参考，使用MySQL来实现过滤。示例：[example/exitscan/dbscan](../../example/exitscan/dbscan)  
>
> **这部分Operator可结合自己的技术进行实现，实现 `ExitFilter` 接口。**



## WithdrawalRequestFilter

代码：

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

`WithdrawalRequestFilter` 用来处理 nETHLargeRequest，因为其也是异步的，需要链下处理。

> 不同Operator可能存在不同的技术实现，我们提供一种最简单示例用以参考，使用MySQL来实现过滤。示例：[example/exitscan/dbscan](../../example/exitscan/dbscan)
>
> **这部分Operator可结合自己的技术进行实现，实现 `WithdrawalRequestFilter` 接口。**



## WithdrawalRequestExitValidatorCounter

代码：

```go
// WithdrawalRequestExitValidatorCounter Calculate the number of validators that need to be exited by a Withdrawal Request
type WithdrawalRequestExitValidatorCounter interface {
	// ExitCounter Calculate the number of validators that need to be exited by a Withdrawal Request
	// @param filterWithdrawalRequests  A list of offline filtered Withdrawal Requests
	ExitCounter(filterWithdrawalRequests []*WithdrawalRequest) (uint32, error)
}
```

实现代码可以参看：[validator/exitscan/neth_large_request.go](../../validator/exitscan/neth_large_request.go)

nETH largeRequest的退出Validator数量，主要根据过滤的 `filterWithdrawalRequests` 列表进行计算。

简要的计算逻辑是根据 计算 `filterWithdrawalRequests`  的 `ClaimEthAmount` 之和，计算退出的validator的数量。

> 计算规则：
>
> if sumETHAmount = 64 ether, need to exit 2 validator
> if sumETHAmount = 66 ether, need to exit 3 validator



## ExitMarker

代码：

```go
// ExitMarker To perform a validator exit, it needs to be flagged, and then it is used for filter
// --------------------------------------------
// The simplest way to implement the operator is to use db, see example
type ExitMarker interface {
	// ExitMark Mark the exit of the Vnft Record
	ExitMark(operatorId *big.Int, records []interface{}) error
}
```

当Validator发起退出后，将其标记为 exited，以便于 Filter 进行过滤。

>  不同Operator可能存在不同的技术实现，我们提供一种最简单示例用以参考，使用MySQL来实现过滤。示例：[example/exitscan/dbscan](../../example/exitscan/dbscan)
>
> **这部分Operator可结合自己的技术进行实现，实现 `ExitMarker` 接口。**



## WithdrawalRequestMarker

`WithdrawalRequest`处理过后也需要标记，以便下一次去filter。

代码：

```go
// WithdrawalRequestMarker To mark deal for WithdrawalRequest.
// --------------------------------------------
// The simplest way to implement the operator is to use db, see example.
type WithdrawalRequestMarker interface {
	// WithdrawalRequestMark To mark deal for WithdrawalRequest
	WithdrawalRequestMark(operatorId *big.Int, withdrawalRequests []*WithdrawalRequest) error
}
```

>  不同Operator可能存在不同的技术实现，我们提供一种最简单示例用以参考，使用MySQL来实现过滤。示例：[example/exitscan/dbscan](../../example/exitscan/dbscan)
>
> **这部分Operator可结合自己的技术进行实现，实现 `WithdrawalRequestMarker` 接口。*





# Example by MySQL

`ExitFilter` 和 `ExitMarker` 的示例使用MySQL。

sql文件参看：[script/sql/operator-sdk-*.sql](../../script/sql/operator-sdk-*.sql)

- `nodedao_validator` 用来存储Validator的必要信息（包括：tokenId、operatorId、是否发起退出、所有者等）
- `neth_withdrawal_request` 用来存储nETH的largeRequest信息，并标记是否处理退出等信息。



## config

一些与配置文件相关的内容，参看：[docs/config.md](../config.md)

