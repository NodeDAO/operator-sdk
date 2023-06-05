// description:
// @author renshiwei
// Date: 2023/6/5

package db

import (
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"strings"
)

// InitMySQL Initialize the MySQL
// @param dsn db conn
// @param logLevel log level
func InitMySQL(dsn, logLevel string) (*gorm.DB, error) {
	if len(dsn) <= 0 {
		return nil, errors.New("dsn is empty")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLog.Default.LogMode(switchGormLog(logLevel)),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to the db.")
	}

	return db, nil
}

func switchGormLog(level string) gormLog.LogLevel {
	var logLevel gormLog.LogLevel
	switch strings.ToLower(level) {
	case "silent":
		logLevel = gormLog.Silent
	case "error":
		logLevel = gormLog.Error
	case "warn":
		logLevel = gormLog.Warn
	case "info":
		logLevel = gormLog.Info
	default:
		logLevel = gormLog.Silent
	}
	return logLevel
}
