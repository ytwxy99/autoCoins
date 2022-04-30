package policy

import (
	"github.com/ytwxy99/autocoins/pkg/interfaces"
	"github.com/ytwxy99/autocoins/pkg/utils"
	"github.com/ytwxy99/autocoins/pkg/utils/index"
)

type MacdPolicy struct{}
type condition struct {
	coin        string
	dataMacd15M []map[string]string
	dataMacd4H  []map[string]string
}

// find macd buy point
func (*MacdPolicy) Target(args ...interface{}) interface{} {
	// convert specified type
	coin := args[0].(string)

	market4H := (&interfaces.MarketArgs{
		CurrencyPair: coin,
		Interval:     -100,
		Level:        utils.Level4Hour,
	}).SpotMarket()
	market15M := (&interfaces.MarketArgs{
		CurrencyPair: coin,
		Interval:     -1,
		Level:        utils.Level15Min,
	}).SpotMarket()

	if market4H != nil && market15M != nil {
		macdArgs := index.DefaultMacdArgs()
		c := &condition{
			coin:        coin,
			dataMacd15M: macdArgs.GetMacd(market15M),
			dataMacd4H:  macdArgs.GetMacd(market4H),
		}

		if len(c.dataMacd4H) < 5 || len(c.dataMacd15M) < 5 {
			return false
		}

		return c.buyCondition()
	}

	return false
}

func (c *condition) buyCondition() bool {
	// judgment depends on 4h price
	conditionA := utils.Compare(c.dataMacd4H[len(c.dataMacd4H)-1][utils.Close], c.dataMacd4H[len(c.dataMacd4H)-1][utils.Open], 0, 0)    //当下4h是具有涨幅的
	conditionB := utils.Compare(c.dataMacd4H[len(c.dataMacd4H)-2][utils.Close], c.dataMacd4H[len(c.dataMacd4H)-3][utils.Close], 0, 1.1) //上个4h是涨幅十个点以上的
	conditionC := utils.Compare(c.dataMacd4H[len(c.dataMacd4H)-3][utils.Close], c.dataMacd4H[len(c.dataMacd4H)-4][utils.Close], 0, 1.05)

	// judgment depends on 4h macd
	conditionD := utils.StringToFloat32(c.dataMacd4H[len(c.dataMacd4H)-1]["macd"]) > 0                                          //当下4h的macd大于0
	conditionE := utils.StringToFloat32(c.dataMacd4H[len(c.dataMacd4H)-2]["macd"]) > 0                                          //上个4h的macd大于0
	conditionF := utils.StringToFloat32(c.dataMacd15M[len(c.dataMacd15M)-1]["macd"]) > 0                                        //当下15m的macd大于0
	conditionG := utils.Compare(c.dataMacd15M[len(c.dataMacd15M)-1]["macd"], c.dataMacd15M[len(c.dataMacd15M)-2]["macd"], 0, 0) //15m的macd是增长的

	// judgment depends on price average data
	averageArgs := index.Average{
		CurrencyPair: c.coin,
		Level:        utils.Level4Hour,
		MA:           utils.Five,
	}
	conditionH := averageArgs.Average(false) > averageArgs.Average(true) //4h的FiveAverage是增长的

	return conditionA && conditionB && !conditionC && conditionD && conditionE && conditionF && conditionG && conditionH
}
