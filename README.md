# operator-sdk

**Using the operator-sdk to integrate with the NodeDAO protocol**.

Made by:   [HashKing](https://www.hashking.com/)




## Config
For more information about the full configuration and its default values, see:  [conf/config-default.yaml](./conf/config-default.yaml)

If you need to use a configuration file to operate, you need to call:
```go
config.InitConfig("<YOUR_CONFIG_FILE.yaml>")
logger.InitLog(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Format)
```
> `config.InitConfig` conf/config-default.yaml is loaded by default, specifying that other configuration files will be merged with configurations.



##  Installation

```go
go get github.com/NodeDAO/operator-sdk
```



## Operator support
Operator integration into Node DAO requires the following:
1. Accumulated over 32 ETH, register Validator. see: [docs/registerValidator](docs/registerValidator)
2. The user initiates unstake, and the operator needs to exit the Validator. see：[docs/exitscan](docs/exitscan)



# License

2023 HashKing

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 2 of the License, or any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the [GNU General Public License](https://github.com/NodeDAO/operator-sdk/blob/main/LICENSE) along with this program. If not, see https://www.gnu.org/licenses/.
