# nETH large Request exitscan 步骤

1、**扫描链上需要退出的nETH largeRequestId 列表**。可通过withdrawalRequest合约`getWithdrawalOfRequestId`方法。

此过程筛选的largeRequestId，还需要进一步过滤，因为异步退出的Validator的状态可能还没有同步到链上。

2、 **在链上获取Operator nETH对应vNFT的记录**。

**nETH对应的vNFT所有者为LiquidStaking，nETH largeRequest 退出的validator必须是 LiquidStaking所拥有的。**

> vNFT方式的记录，所有者为用户。

3、**Operator自行实现已发起退出的Validator过滤**。

Validator的退出需要经过beacon的生命周期，退出是异步的，退出完成时间并不确定。发起退出的Validator在Oracle没有上报触发结算，链上数据是并不会更新。因此需要链下的手段来进行过滤。

不同Operator可能存在不同的技术实现，我们提供一种最简单示例用以参考，使用MySQL来实现过滤。

4、**退出Validator并在数据库中进行标记**。



# 具体实现

## 1、扫描链上需要退出的nETH largeRequestId 列表

实现参看：[validator/exitscan/neth_large_request_test.go](../../validator/exitscan/neth_large_request_test.go)

测试用例参看：[validator/exitscan/neth_large_request_test.go](../../validator/exitscan/neth_large_request_test.go)

方法参看：`TestNethExitScan_ByLocal`

关键代码：

```go
// 初始化 neth exit scan 对象
nethExitScan, err := NewNETHExitScan(ctx, config.GlobalConfig.Eth.Network, config.GlobalConfig.Eth.ElAddr)
require.NoError(t, err)

// 扫描智能合约中需要处理的 WithdrawalRequest
withdrawalInfo, err := nethExitScan.WithdrawalRequestScan(big.NewInt(1))
require.NoError(t, err)
fmt.Printf("withdrawalInfo scan:%+v\n", withdrawalInfo)

// 扫描智能合约中 未退出并且是 LiquidStaking 所有的vNFT
nethRecords, err := nethExitScan.ExitScan(big.NewInt(1))
require.NoError(t, err)
```



## 2、

