# Register Validator on NodeDAO

See Solidity source code：https://github.com/NodeDAO/NodeDAO-Protocol/blob/main/src/LiquidStaking.sol

The user stakes ETH to the staking pool of the NodeDAO LiquidStaking corresponding to the operator, the operator should be responsible for `registerValidator` to the NodeDAO protocol, and the `registerValidator` will `deposit` 32ETH.


# Implementation steps

## 1. Monitor whether the pool is greater than 32 ETH

Scan the LiquidStaking contract regularly to determine whether the Operator's pool accumulation is greater than 32ETH. To accumulate more than 32 ETH, a corresponding number of validators is required for `registerValidator`.
```go
registerValidatorCount = operatorPoolBalances/32
```

The way to monitor operatorPoolBalances is in the LiquidStaking contract:
```solidity
// operator's internal stake pool, key is operator_id
mapping(uint256 => uint256) public operatorPoolBalances;
```



## 2. Generate keystore and deposit data

Based on the registerValidatorCount calculated in the previous step, a fixed number of depositData and keystores are generated.

> generate keystore+password, or you can generate a private key directly in order to run Validator.
>
> depositData is to stake 32ETH to beacon.

Tool Recommendations:

- [ethereum/staking-deposit-cli](https://github.com/ethereum/staking-deposit-cli) Keystore and deposit data files can be generated, official tools.
- [ethdo depositData](https://github.com/wealdtech/ethdo/tree/master/cmd/validator/depositdata)
- [wealdtech/go-eth2-wallet-encryptor-keystorev4](https://github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4)



## 3、NodeDAO registerValidator

> NodeDAO Operator needs to set the controllerAddress when registering, and you need to use the controllerAddress to initiate the `registerValidator`. And make sure controllerAddress has enough ETH to pay for gas.
>
> `registerValidator` will register the store contract with the eth of operatorPoolBalances.

For code implementation, see: [validator/register/register_validator.go](../../validator/register/register_validator.go)

## 4. Start Validator

According to the keystore+password generated in step 3, start the Validator program.

> This part can be implemented according to the operator's technology stack.

See also:

- [teku](https://docs.teku.consensys.net/)
- [prysm](https://docs.prylabs.network/docs/getting-started)
- ......


## 5、Other: for exit scan example

In the exitscan example, we provide an example implementation of mysql for Validator's one-step exit filtering. `exitscan filter`

You need to rely on registerValidator data.

After starting Validator, you can store Validator's data in a mysql table: `nodedao_validator`.

> This part, the operator can be implemented according to the specific situation.

See:

- sql：[script/sql/operator-sdk-*.sql](../../script/sql/operator-sdk-*.sql)
- Example: [example/registerdb/register_validator_example.go](../../example/registerdb/register_validator_example.go)


## Scheduled tasks are executed periodically

The above three steps should be performed periodically in a timed task, a user initiates unstake, and the operator should exit the Validator in time.

> Timed task implementation: omitted