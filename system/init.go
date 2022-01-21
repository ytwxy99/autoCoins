package system

import (
	"github.com/gateio/gateapi-go/v6"
	"github.com/ytwxy99/autoCoins/client"
	"github.com/ytwxy99/autoCoins/configuration"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
)

// get all usdt pair coins and write into file
func InitCurrencyPairs(client *gateapi.APIClient, pairs []gateapi.CurrencyPair, filePath string) error {
	coins := []string{}
	for _, pair := range pairs {
		// just record coin which is tradable
		if pair.TradeStatus == "tradable" && pair.Quote == "USDT" {
			values := interfaces.K(client, pair.Id, -1, "1d")
			if utils.StringToFloat32(values[0][1]) >= 200000.0 {
				coins = append(coins, pair.Id)
			}
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
