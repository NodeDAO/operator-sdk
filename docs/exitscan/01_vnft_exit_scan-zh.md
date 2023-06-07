# vNFT exitscan 步骤

1、**扫描链上需要退出的vNFT的tokenId**，可通过withdrawalRequest合约`getUserUnstakeButOperatorNoExitNfs`方法。

此过程筛选的tokenId，还需要进一步过滤，因为异步退出的Validator的状态可能还没有同步到链上。

2、**Operator自行实现已发起退出的Validator过滤**。

Validator的退出需要经过beacon的生命周期，退出是异步的，退出完成时间并不确定。发起退出的Validator在Oracle没有上报触发结算，链上数据是并不会更新。因此需要链下的手段来进行过滤。

不同Operator可能存在不同的技术实现，我们提供一种最简单示例用以参考，使用MySQL来实现过滤。

3、**退出Validator并在数据库中进行标记**。


# 具体实现

## 1、扫描链上待退出的vNFT

实现参看：[validator/exitscan/vnft.go](../../validator/exitscan/vnft.go)

测试用例参看：[validator/exitscan/vnft_test.go](../../validator/exitscan/vnft_test.go)   

方法：`TestVnftExitScan_ByLocal`

关键代码：

```go
// 初始化 vnft exit scan 对象
vnftExitScan, err := NewVnftExitScan(ctx, config.GlobalConfig.Eth.Network, config.GlobalConfig.Eth.ElAddr)
require.NoError(t, err)

// 扫描智能合约中需要退出的 vnftRecords
vnftRecords, err := vnftExitScan.ExitScan(big.NewInt(1))
require.NoError(t, err)
```



## 2、过滤已发起退出的Validator

> 此部分示例 这里借助MySQL实现，Operator可以根据具体情况进行实现。
>
> 在过滤退出的Validator前，Validator应当已经在数据库表 `nodedao_validator` 中存在。可以再 `registerValidator`后将其存储到库中，这部分会在 `registerValidator` 操作里说明。现在可以假设 Validator 已经在表 `nodedao_validator` 中。

代码实现参看：[example/exitscan/dbscan/vnft.go](../../example/exitscan/dbscan/vnft.go) 

方法：`Filter`

`vnftExitScan.ExitScan` 在链上扫描的记录 `vnftRecords` ，作为 `Filter` 的入参 `vnftContractExitRecords`。



## 3、退出Validator并在数据库中进行标记

> 此部分示例 这里借助MySQL实现，Operator可以根据具体情况进行实现。

第二步中过滤后的Validators即是Operator需要退出的，Operator需要根据自己的技术实现去发起Validator的退出。 并标记这些Validator已经发起退出，用于下一次的 `Filter`。

代码实现参看：[example/exitscan/dbscan/vnft.go](../../example/exitscan/dbscan/vnft.go) 

方法：`ExitMark`



## 定时任务 周期性执行

上述三个步骤应当在 **定时任务** 中周期执行，有用户发起 unstake，Operator应当及时退出Validator。

定时任务实现：略