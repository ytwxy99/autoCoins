package policy

import (
	"gorm.io/gorm"

	"github.com/ytwxy99/autocoins/pkg/interfaces"
	"github.com/ytwxy99/autocoins/pkg/utils"
	"github.com/ytwxy99/autocoins/pkg/utils/index"
)

type Trend30M struct{}

// cointegration policy
func (Trend30M) Target(args ...interface{}) interface{} {
	isBuy := map[string]string{}
	db := args[0].(*gorm.DB)
	//sysConf := args[1].(*configuration.SystemConf)
	coin := args[2].(string)

	sports := (&interfaces.MarketArgs{
		CurrencyPair: coin,
		Interval:     utils.Now,
		Level:        utils.Level30Min,
	}).SpotMarket()
	if sports == nil {
		return isBuy
	}

	currentPrice := utils.StringToFloat64(sports[0][2])
	btcRisingCondition := conditionUpMonitor30M(coin, 1.001, currentPrice)
	btcFallingCondition := conditionDownMonitor30M(coin, 1.001, currentPrice)

	// rising market buy point
	if btcRisingCondition {
		if tradeJugde(coin, db, "up") {
			isBuy[coin] = utils.DirectionUp
			return isBuy
		}

	}

	// falling market buy point
	if btcFallingCondition {
		if tradeJugde(coin, db, "down") {
			isBuy[coin] = utils.DirectionDown
			return isBuy
		}
	}

	return isBuy
}

func conditionUpMonitor30M(coin string, averageDiff float64, currentPrice float64) bool {
	averageArgs := &index.Average{
		CurrencyPair: coin,
		Level:        utils.Level30Min,
		MA:           utils.MA21,
	}
	MA21Average := averageArgs.Average(false) > averageArgs.Average(true)

	averageArgs.MA = utils.MA5
	MA5Average := averageArgs.Average(false) > averageArgs.Average(true)*averageDiff

	priceC := currentPrice > averageArgs.Average(false)

	return MA21Average && MA5Average && priceC
}

func conditionDownMonitor30M(coin string, averageDiff float64, currentPrice float64) bool {
	averageArgs := &index.Average{
		CurrencyPair: coin,
		Level:        utils.Level30Min,
		MA:           utils.MA21,
	}
	MA21Average := averageArgs.Average(false)*averageDiff < averageArgs.Average(true) //4h的Average是增长的

	averageArgs.MA = utils.MA5
	MA5Average := averageArgs.Average(false)*averageDiff < averageArgs.Average(true)

	priceC := currentPrice < averageArgs.Average(false)

	return MA21Average && MA5Average && priceC
}
