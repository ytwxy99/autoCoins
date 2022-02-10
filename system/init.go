package system

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/client"
	"github.com/ytwxy99/autoCoins/configuration"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
)

// get all usdt pair coins and write into file
func InitCurrencyPairs(client *gateapi.APIClient, pairs []gateapi.CurrencyPair, filePath string, db *gorm.DB) error {
	var historyDay *database.HistoryDay
	coins := []string{}
	for _, pair := range pairs {
		// just record coin which is tradable
		if pair.TradeStatus == "tradable" && pair.Quote == "USDT" {
			values := interfaces.K(client, pair.Id, -365, "1d")
			for i, value := range values {
				fmt.Println(i, value[0], value[2], pair.Id)
				timeTrans, err := time.Parse(" 2022-02-10 08:00:00", value[0])
				if err != nil {
					logrus.Error("get time type from string error: %s\n", err)
				}

				historyDay = &database.HistoryDay{
					Contract: pair.Id,
					Time:     timeTrans,
					Price:    value[2],
				}

				err = historyDay.AddHistoryDay(db)
				if err != nil {
					logrus.Error("add HistoryDay record error: %s\n", err)
				}
			}
			if utils.StringToFloat32(values[0][1]) >= 200000.0 {
				coins = append(coins, pair.Id)
			}
			break
		}
	}

	return utils.WriteLines(coins, filePath)
}

// init system base data
func Init(authConf *configuration.GateAPIV4, sysConf *configuration.SystemConf) {
	c, ctx := client.GetSpotClient(authConf)
	utils.InitLog(sysConf.LogPath)
	db := database.GetDB(sysConf)
	database.InitDB(db)
	InitCmd(c, ctx, sysConf, db)
}
