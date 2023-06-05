// description:
// @author renshiwei
// Date: 2022/10/5 17:06

//go:generate go-bindata -pkg=config -nocompress -o=default_conf.go ../conf/config-default.yaml

package config

import (
	"github.com/NodeDAO/operator-sdk/common/db"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once

var (
	GlobalConfig Config
	GlobalDB     *gorm.DB
)

type Config struct {
	Cli struct {
		Name string
	}

	Log struct {
		Level  string
		Format string
	}

	DB struct {
		Dsn      string
		LogLevel string
	}

	Eth struct {
		Network    string
		ElAddr     string
		ClAddr     string
		PrivateKey string
	}
}

// InitOnce The operation is initialized only once.
// Before, you need to call 'config.InitConfig`.
func InitOnce() error {
	var err error

	once.Do(func() {
		GlobalDB, err = db.InitMySQL(GlobalConfig.DB.Dsn, GlobalConfig.DB.LogLevel)
	})

	if err != nil {
		return errors.Wrap(err, "Failed to InitOnce.")
	}

	return nil
}
