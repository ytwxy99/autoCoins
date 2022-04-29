package trade

import (
	"math"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/pkg/configuration"
	"github.com/ytwxy99/autoCoins/pkg/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
	"github.com/ytwxy99/autoCoins/utils/index"
)

type SellArgs struct {
	Policy         string
	Contract       string
	db             *gorm.DB
	LastPrice      float32
	StoredPrice    float32
	OrderDirection string
	sysConfig      *configuration.SystemConf
}

func (sellArgs *SellArgs) SellPolicy() bool {
	if sellArgs.Policy == "trend" {
		// sell specify coin when the absolute value of rising or falling rate over 15%.
		return math.Abs(float64((sellArgs.LastPrice-sellArgs.StoredPrice)/sellArgs.StoredPrice)) >= 15
	} else if sellArgs.Policy == "cointegration" {
		weightResult := map[string]bool{}
		order, err := (database.Order{
			Contract:  sellArgs.Contract,
			Direction: sellArgs.OrderDirection,
		}).FetchOneOrder(sellArgs.db)
		if err != nil {
			logrus.Error("Can't find coin in Order record, then trade will be canceled : ", sellArgs.Contract)
			return false
		}

		if order.Direction == utils.DirectionUp {
			// sell it when the coin has falled over 3%
			sports := (&interfaces.MarketArgs{
				CurrencyPair: sellArgs.Contract,
				Interval:     utils.Now,
				Level:        utils.Level4Hour,
			}).SpotMarket()
			currentPrice := utils.StringToFloat64(sports[0][2])
			if (currentPrice-utils.StringToFloat64(order.Price))/utils.StringToFloat64(order.Price) < -0.03 {
				return true
			}

			// sell it when 21 MA is falling
			average := &index.Average{
				CurrencyPair: sellArgs.Contract,
				Level:        utils.Level4Hour,
				MA:           utils.MA21,
			}
			M1 := average.Average(false)
			M2 := average.Average(true)
			if M1 < M2 && M1 != 0 && M2 != 0 {
				return true
			}

			weights, err := utils.ReadLines(sellArgs.sysConfig.WeightCsv)
			if err != nil {
				logrus.Error("read weight coins failed, the error is ", weights)
			}
			for _, weight := range weights {
				average.CurrencyPair = weight
				weightMa21 := average.Average(false)
				weightSport := (&interfaces.MarketArgs{
					CurrencyPair: weight,
					Interval:     utils.Now,
					Level:        utils.Level4Hour,
				}).SpotMarket()

				if utils.StringToFloat64(weightSport[0][2]) < weightMa21 {
					weightResult[weight] = true
				} else {
					weightResult[weight] = false
				}
			}

			if currentPrice < M1 {
				count := 0
				for _, value := range weightResult {
					if value {
						count++
					}
				}

				if float64(count)/float64(len(weights)) > 0.7 {
					return true
				}
			}
			return false
		} else {
			// sell it when the coin has rised over 3%
			sports := (&interfaces.MarketArgs{
				CurrencyPair: sellArgs.Contract,
				Interval:     utils.Now,
				Level:        utils.Level4Hour,
			}).SpotMarket()
			currentPrice := utils.StringToFloat64(sports[0][2])
			if (currentPrice-utils.StringToFloat64(order.Price))/utils.StringToFloat64(order.Price) > 0.03 {
				return true
			}

			// sell it when 21 MA is falling
			average := &index.Average{
				CurrencyPair: sellArgs.Contract,
				Level:        utils.Level4Hour,
				MA:           utils.MA21,
			}
			M1 := average.Average(false)
			M2 := average.Average(true)
			if M1 > M2 && M1 != 0 && M2 != 0 {
				return true
			}

			weights, err := utils.ReadLines(sellArgs.sysConfig.WeightCsv)
			if err != nil {
				logrus.Error("read weight coins failed, the error is ", weights)
			}
			for _, weight := range weights {
				average.CurrencyPair = weight
				weightMa21 := average.Average(false)
				weightSport := (&interfaces.MarketArgs{
					CurrencyPair: weight,
					Interval:     utils.Now,
					Level:        utils.Level4Hour,
				}).SpotMarket()

				if utils.StringToFloat64(weightSport[0][2]) > weightMa21 {
					weightResult[weight] = true
				} else {
					weightResult[weight] = false
				}
			}

			if currentPrice > M1 {
				count := 0
				for _, value := range weightResult {
					if value {
						count++
					}
				}

				if float64(count)/float64(len(weights)) > 0.7 {
					return true
				}
			}
			return false
		}
	} else {
		return false
	}
}
