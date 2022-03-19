package trade

import (
	"math"

	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	py "github.com/ytwxy99/autoCoins/policy"
	"github.com/ytwxy99/autoCoins/utils"
)

func SellPolicy(policy string, args ...interface{}) bool {
	if policy == "macd" {
		lastPrice := args[0].(float32)
		storedPrice := args[1].(float32)

		return math.Abs(float64((lastPrice-storedPrice)/storedPrice)) >= 15
	} else if policy == "cointegration" {
		//lastPrice := args[0].(float32)
		//storedPrice := args[1].(float32)
		coin := args[2].(string)
		client := args[3].(*gateapi.APIClient)
		db := args[5].(*gorm.DB)

		tradeDetail := database.TradeDetail{
			Contract: coin,
		}

		record, err := tradeDetail.FetchOneTradeDetail(db)
		if err != nil {
			logrus.Info("Can't find coin-pair in inOrder record, then trade will be canceled : ", coin)
			return false
		}

		// conditions no.1
		//if math.Abs(float64((lastPrice-storedPrice)/storedPrice)) >= 15 {
		//	return true
		//}

		// conditions no.2
		k15mValues := interfaces.K(client, record.CointPair, -10, "15m")
		if k15mValues != nil {
			k4mMacds := py.GetMacd(k15mValues, 12, 26, 9)
			if len(k4mMacds) < 10 {
				return false
			}

			macdValue := utils.StringToFloat32(k4mMacds[len(k4mMacds)-1]["macd"])
			if macdValue < 0 {
				logrus.Info("Find sell point, then sell it : ", record.Contract, "Macd-Value", macdValue)
				return true
			}
		}

		return false
	} else {
		return false
	}
}
