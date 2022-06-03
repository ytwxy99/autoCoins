package trade

import (
	"math"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ytwxy99/autocoins/database"
	"github.com/ytwxy99/autocoins/pkg/configuration"
	"github.com/ytwxy99/autocoins/pkg/interfaces"
	"github.com/ytwxy99/autocoins/pkg/utils"
	"github.com/ytwxy99/autocoins/pkg/utils/index"
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
	} else if sellArgs.Policy == "trend30m" {
		average := &index.Average{
			CurrencyPair: sellArgs.Contract,
			Level:        utils.Level30Min,
			MA:           utils.MA21,
		}

		sports := (&interfaces.MarketArgs{
			CurrencyPair: average.CurrencyPair,
			Interval:     utils.Now,
			Level:        utils.Level30Min,
		}).SpotMarket()
		currentPrice := utils.StringToFloat64(sports[0][2])

		order, err := (&database.Order{
			Contract:  sellArgs.Contract,
			Direction: sellArgs.OrderDirection,
		}).FetchOneOrder(sellArgs.db)
		if err != nil {
			logrus.Error("fetch orders failed: ", err)
			return false
		}

		// for rising
		if sellArgs.OrderDirection == utils.DirectionUp {
			r1 := average.Average(false) <= average.Average(true)
			r2 := order != nil && currentPrice*1.01 <= utils.StringToFloat64(order.Price)
			r3 := currentPrice < average.Average(false)

			average.MA = utils.MA5
			r4 := average.Average(false) <= average.Average(true)

			if r1 || r2 || r3 || r4 {
				// do sell
				if order != nil {
					sold := &database.Sold{
						Contract:        sellArgs.Contract,
						Price:           utils.Float32ToString(float32(currentPrice)),
						Volume:          0,
						Time:            utils.GetNowTimeStamp(),
						Profit:          float32(currentPrice) - utils.StringToFloat32(order.Price),
						Relative_profit: utils.Float32ToString((float32(currentPrice) - utils.StringToFloat32(order.Price)) / utils.StringToFloat32(order.Price) * 100),
						Test:            "test-order",
						Status:          "close",
						Typee:           "limit",
						Account:         "spot",
						Side:            "sell",
						Iceberg:         "0",
					}
					if err := sold.AddSold(sellArgs.db); err != nil {
						logrus.Errorf("Craete sold failed, the sold detail is %s, the error is %s.", sold, err)
						return false
					}
				}

				return true
			}
		} else if sellArgs.OrderDirection == utils.DirectionDown {
			// for falling
			average.MA = utils.MA21
			f1 := average.Average(false) >= average.Average(true)
			f2 := order != nil && currentPrice >= utils.StringToFloat64(order.Price)*1.01
			f3 := currentPrice > average.Average(false)

			average.MA = utils.MA5
			f4 := average.Average(false) >= average.Average(true)

			if f1 || f2 || f3 || f4 {
				// do sell
				if order != nil {
					sold := &database.Sold{
						Contract:        sellArgs.Contract,
						Price:           utils.Float32ToString(float32(currentPrice)),
						Volume:          0,
						Time:            utils.GetNowTimeStamp(),
						Profit:          float32(currentPrice) - utils.StringToFloat32(order.Price),
						Relative_profit: utils.Float32ToString((float32(currentPrice) - utils.StringToFloat32(order.Price)) / utils.StringToFloat32(order.Price) * 100),
						Test:            "test-order",
						Status:          "close",
						Typee:           "limit",
						Account:         "spot",
						Side:            "sell",
						Iceberg:         "0",
					}
					if err := sold.AddSold(sellArgs.db); err != nil {
						logrus.Errorf("Craete sold failed, the sold detail is %s, the error is %s.", sold, err)
						return false
					}
				}

				return true
			}
		}
	}

	return false
}
