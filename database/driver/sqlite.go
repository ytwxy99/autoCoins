package driver

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ytwxy99/autocoins/pkg/configuration"
)

type SqliteDrive struct {
	Conf *configuration.SystemConf
}

func (sqliteDrive *SqliteDrive) DatabaseConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqliteDrive.Conf.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Open database error!")
	}

	database, err := db.DB()
	if err != nil {
		panic("connect db server failed.")
	}

	database.SetMaxIdleConns(10)                   // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	database.SetMaxOpenConns(100)                  // SetMaxOpenConns sets the maximum number of open connections to the database.
	database.SetConnMaxLifetime(time.Second * 600) // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.

	return db
}
