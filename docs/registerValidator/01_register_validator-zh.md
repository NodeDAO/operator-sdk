# Register Validator on NodeDAO

solidity源代码参看：https://github.com/NodeDAO/NodeDAO-Protocol/blob/main/src/LiquidStaking.sol

用户将 ETH 质押到 NodeDAO LiquidStaking 对应 Operator 的质押池当中，Operator应当负责 `registerValidator` 到 NodeDAO协议当中，`registerValidator`  会进行 32ETH 的 `deposit`。



# 实现步骤

## 1、监控pool是否大于32ETH

定时扫描LiquidStaking合约，判断Operator的pool是否累积大于32ETH。累积超过32ETH，需要 `registerValidator` 相应数量的Validator。

```go
registerValidatorCount = operatorPoolBalances/32
```

监控operatorPoolBalances的方法在 LiquidStaking 合约：

```solidity
// operator's internal stake pool, key is operator_id
mapping(uint256 => uint256) public operatorPoolBalances;
```



## 2、生成keystore和depositData

根据上一步计算的registerValidatorCount，生成固定数量的 depositData和keystore。

> 生成keystore+password，或者可以直接生成私钥，目的是为了运行Validator。
>
> depositData是为了将32ETH质押到beacon。

工具推荐：

- [ethereum/staking-deposit-cli](https://github.com/ethereum/staking-deposit-cli) 可生成keystore和depositData文件，官方工具。
- [ethdo depositData](https://github.com/wealdtech/ethdo/tree/master/cmd/validator/depositdata)
- [wealdtech/go-eth2-wallet-encryptor-keystorev4](https://github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4) 



## 3、NodeDAO registerValidator

> NodeDAO Operator注册时需要设置 controllerAddress，需要使用controllerAddress用来发起 `registerValidator` 。并确保controllerAddress有足够的ETH来支付gas。
>
> `registerValidator` 会将operatorPoolBalances的eth注册的deposit合约。

代码实现参看：[validator/register/register_validator.go](../../validator/register/register_validator.go)



## 4、启动Validator

根据第三步生成的keystore+password，去启动Validator程序。

> 这一部分可以根据Operator的技术栈去实现。

可参看：

- [teku](https://docs.teku.consensys.net/)
- [prysm](https://docs.prylabs.network/docs/getting-started)
- ……



## 5、Other: for exit scan example

在 exitscan 的的例子中，对于Validator的一步退出过滤，我们提供了mysql的实现示例。``exitscan filter`

需要依赖registerValidator的数据。

启动Validator后，可以将Validator的数据存储到mysql的表：`nodedao_validator`。

> 这一部分，Operator可根据具体情况实现。

参看：

- sql：[script/sql/operator-sdk-*.sql](../../script/sql/operator-sdk-*.sql)
- Example: [example/registerdb/register_validator_example.go](../../example/registerdb/register_validator_example.go)



## 定时任务 周期性执行

上述三个步骤应当在 **定时任务** 中周期执行，有用户发起 unstake，Operator应当及时退出Validator。

> 定时任务实现：略