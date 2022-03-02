package system

import (
	"fmt"
	"gorm.io/gorm"
	"os/exec"
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
			values := interfaces.K(client, pair.Id, -730, "1d")
			for _, value := range values {
				timeTrans, err := time.ParseInLocation("2006-01-02 08:00:00", value[0], time.Local)
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
					if err.Error() != utils.DBHistoryDayUniq {
						logrus.Error("add HistoryDay record error: %s\n", err)
					}
				}
			}

			if len(values) != 0 && utils.StringToFloat32(values[0][1]) >= 200000.0 {
				coins = append(coins, pair.Id)
			}
		}
	}

	return utils.WriteLines(coins, filePath)
}

func InitCointegration(dbPath string, scriptPath string, coinCsv string) error {
	cmd := exec.Command("python3", scriptPath, dbPath, coinCsv)
	output, err := cmd.Output()
	fmt.Println(string(output))
	if err != nil {
		logrus.Error("run cointegration python srcipt error:", err)
		return err
	}

	return nil
}

// init system base data
func Init(authConf *configuration.GateAPIV4, sysConf *configuration.SystemConf) {
	c, ctx := client.GetSpotClient(authConf)
	utils.InitLog(sysConf.LogPath)
	db := database.GetDB(sysConf)
	database.InitDB(db)
	InitCmd(c, ctx, sysConf, db)
}
