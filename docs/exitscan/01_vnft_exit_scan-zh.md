# vNFT exitscan 步骤

1、**扫描链上需要退出的vNFT的tokenId**，可通过withdrawalRequest合约`getUserUnstakeButOperatorNoExitNfs`方法。

此过程筛选的tokenId，还需要进一步过滤，因为异步退出的Validator的状态可能还没有同步到链上。

2、**Operator自行实现已发起退出的Validator过滤**。

Validator的退出需要经过beacon的生命周期，退出是异步的，退出完成时间并不确定。发起退出的Validator在Oracle没有上报触发结算，链上数据是并不会更新。因此需要链下的手段来进行过滤。

不同Operator可能存在不同的技术实现，我们提供一种最简单示例用以参考，使用MySQL来实现过滤。



# 具体实现

## 1、扫描链上待退出的vNFT

实现参看：[validator/exitscan/vnft.go](../../validator/exitscan/vnft.go)

测试用例参看：[validator/exitscan/vnft_test.go](../../validator/exitscan/vnft_test.go)   方法：`TestVnftExitScan_ByLocal`

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

