
## 生成abi
```shell
solc --abi --bin <your-contract.sol> -o <output-directory>
```

## 生成合约代码

```shell
abigen --abi=contracts/liq/LiquidStaking.json --pkg=liq --out=contracts/liq/LiquidStaking.go

abigen --abi=contracts/vnft/VNFT.json --pkg=vnft --out=contracts/vnft/vnft.go

abigen --abi=contracts/withdrawalRequest/withdrawalRequest.json --pkg=withdrawalRequest --out=contracts/withdrawalRequest/withdrawal_request.go
```
