package trade

import (
	"math"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
	"github.com/ytwxy99/autoCoins/utils/index"
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
		db := args[4].(*gorm.DB)

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
		k15mValues := (&interfaces.MarketArgs{
			CurrencyPair: record.CointPair,
			Interval:     -10,
			Level:        utils.Level15Min,
		}).Market()
		if k15mValues != nil {
			macdArgs := index.DefaultMacdArgs()
			k15mMacds := macdArgs.GetMacd(k15mValues)
			if len(k15mMacds) < 10 {
				return false
			}

			macd15mValue := utils.StringToFloat32(k15mMacds[len(k15mMacds)-1]["macd"])
			if macd15mValue < 0 {
				logrus.Info("Find sell point, then sell it : ", record.Contract, "Macd-Value", macd15mValue)
				return true
			}
		}

		return false
	} else {
		return false
	}
}
