package driver

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ytwxy99/autoCoins/pkg/configuration"
)

type MysqlDrive struct {
	Conf *configuration.SystemConf
}

func (mysqlDrive *MysqlDrive) DatabaseConnect() *gorm.DB {
	conf := mysqlDrive.Conf.Mysql
	mysqlUrl := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	db, err := gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Open database error!")
	}

	return db
}
