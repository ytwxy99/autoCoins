package policy

import (
	"context"
	"github.com/ytwxy99/autocoins/pkg/interfaces"
	"github.com/ytwxy99/autocoins/pkg/utils"
	"github.com/ytwxy99/autocoins/pkg/utils/index"
)

type TrendPolicy struct {
}

// Target find macd buy point
func (*TrendPolicy) Target(ctx context.Context) map[string]string {
	isBuy := make(map[string]string)
	coin := ctx.Value("coin").(string)

	sports := (&interfaces.MarketArgs{
		CurrencyPair: coin,
		Interval:     utils.Now,
		Level:        utils.Level30Min,
	}).SpotMarket()
	if sports == nil {
		return isBuy
	}

	//currentPrice := utils.StringToFloat64(sports[0][2])
	if trendRising(coin) {
		// for rising market
		isBuy[coin] = utils.DirectionUp
	}

	if trendfalling(coin) {
		// for falling market
		isBuy[coin] = utils.DirectionDown
	}

	return isBuy
}

func trendRising(coin string) bool {
	// 均线判断
	averageArgs := &index.Average{
		CurrencyPair: coin,
		Level:        utils.Level4Hour,
		MA:           utils.MA21,
	}
	MA21 := averageArgs.Average(false) > averageArgs.Average(true)
	averageArgs.MA = utils.MA10
	MA10 := averageArgs.Average(false) > averageArgs.Average(true)
	averageArgs.MA = utils.MA5
	MA5 := averageArgs.Average(false) > averageArgs.Average(true)

	// macd判断
	if MA21 && MA10 && MA5 && conditionRising(coin) {
		return true
	} else {
		return false
	}
}

func trendfalling(coin string) bool {
	// 均线判断
	averageArgs := &index.Average{
		CurrencyPair: coin,
		Level:        utils.Level4Hour,
		MA:           utils.MA21,
	}
	MA21 := averageArgs.Average(false) < averageArgs.Average(true)
	averageArgs.MA = utils.MA10
	MA10 := averageArgs.Average(false) < averageArgs.Average(true)
	averageArgs.MA = utils.MA5
	MA5 := averageArgs.Average(false) < averageArgs.Average(true)

	// macd判断
	if MA21 && MA10 && MA5 && conditionFalling(coin) {
		return true
	} else {
		return false
	}
}

func conditionRising(contract string) bool {

	dataMacd4h := index.GetMacd(contract, utils.Level4Hour)
	dataMacd15m := index.GetMacd(contract, utils.Level4Hour)
	if len(dataMacd4h) < 5 || len(dataMacd15m) < 5 {
		return false
	}

	// judgment depends on 4h macd
	a := utils.StringToFloat32(dataMacd4h[len(dataMacd4h)-1]["macd"]) > 0 //当下4h的macd大于0
	//conditionE := utils.StringToFloat32(c.dataMacd4H[len(c.dataMacd4H)-2]["macd"]) > 0   //上个4h的macd大于0
	b := utils.StringToFloat32(dataMacd15m[len(dataMacd15m)-1]["macd"]) > 0                                    //当下15m的macd大于0
	c := utils.Compare(dataMacd15m[len(dataMacd15m)-1]["macd"], dataMacd15m[len(dataMacd15m)-2]["macd"], 0, 0) //15m的macd是增长的

	return a && b && c
}

func conditionFalling(contract string) bool {

	dataMacd4h := index.GetMacd(contract, utils.Level4Hour)
	dataMacd15m := index.GetMacd(contract, utils.Level4Hour)
	if len(dataMacd4h) < 5 || len(dataMacd15m) < 5 {
		return false
	}

	// judgment depends on 4h macd
	a := utils.StringToFloat32(dataMacd4h[len(dataMacd4h)-1]["macd"]) < 0 //当下4h的macd小于0
	//conditionE := utils.StringToFloat32(c.dataMacd4H[len(c.dataMacd4H)-2]["macd"]) < 0   //上个4h的macd小于0
	b := utils.StringToFloat32(dataMacd15m[len(dataMacd15m)-1]["macd"]) < 0                                    //当下15m的macd小于0
	c := utils.Compare(dataMacd15m[len(dataMacd15m)-1]["macd"], dataMacd15m[len(dataMacd15m)-2]["macd"], 0, 0) //15m的macd是增长的

	return a && b && !c
}
