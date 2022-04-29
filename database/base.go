package database

import (
	"gorm.io/gorm"

	"github.com/ytwxy99/autoCoins/database/driver"
	"github.com/ytwxy99/autoCoins/pkg/configuration"
)

type DatabaseDrive interface {
	DatabaseConnect() *gorm.DB
}

func GetDB(sysConf *configuration.SystemConf) *gorm.DB {
	var databaseDrive DatabaseDrive

	if sysConf.DBType == "sqlite" {
		databaseDrive = &driver.SqliteDrive{
			Conf: sysConf,
		}
		return databaseDrive.DatabaseConnect()

	} else if sysConf.DBType == "mysql" {
		databaseDrive = &driver.MysqlDrive{
			Conf: sysConf,
		}
		return databaseDrive.DatabaseConnect()
	}

	return nil
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
