package policy

import (
	"context"

	"github.com/ytwxy99/autocoins/pkg/interfaces"
	"github.com/ytwxy99/autocoins/pkg/utils"
	"github.com/ytwxy99/autocoins/pkg/utils/index"
)

type Trend30M struct{}

// Target cointegration policy
func (Trend30M) Target(ctx context.Context) interface{} {
	isBuy := map[string]string{}
	coin := ctx.Value("coin").(string)

	sports := (&interfaces.MarketArgs{
		CurrencyPair: coin,
		Interval:     utils.Now,
		Level:        utils.Level30Min,
	}).SpotMarket()
	if sports == nil {
		return isBuy
	}

	currentPrice := utils.StringToFloat64(sports[0][2])
	risingCondition := conditionUpMonitor30M(coin, 1.001, currentPrice)
	fallingCondition := conditionDownMonitor30M(coin, 1.001, currentPrice)

	// rising market buy point
	if risingCondition {
		if tradeJugde(ctx, coin, "up") && utils.IsTradeTime() {
			isBuy[coin] = utils.DirectionUp
			return isBuy
		}

	}

	// falling market buy point
	if fallingCondition {
		if tradeJugde(ctx, coin, "down") && utils.IsTradeTime() {
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
	// the 21th average line is rising
	MA21Average := averageArgs.Average(false) > averageArgs.Average(true)

	// the 5th average line is rising
	averageArgs.MA = utils.MA5
	MA5Average := averageArgs.Average(false) > averageArgs.Average(true)*averageDiff

	// the pricing is rising and larger than 21th average line
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
