package database

import (
	"time"

	"github.com/ytwxy99/autoCoins/configuration"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB(sysConf *configuration.SystemConf) (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(sysConf.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Open database error!")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("connect db server failed.")
	}

	sqlDB.SetMaxIdleConns(10)                   // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxOpenConns(100)                  // SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetConnMaxLifetime(time.Second * 600) // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.

	return db
}

func InitDB(db *gorm.DB) {
	// init database
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Sold{})
	db.AutoMigrate(&InOrder{})
	db.AutoMigrate(&HistoryDay{})
	db.AutoMigrate(&Cointegration{})
	db.AutoMigrate(&TradeDetail{})
}
