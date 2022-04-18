package policy

import (
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ytwxy99/autoCoins/configuration"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
	"github.com/ytwxy99/autoCoins/utils/index"
)

type Cointegration struct{}

// cointegration policy
func (Cointegration) Target(args ...interface{}) interface{} {
	db := args[0].(*gorm.DB)
	sysConf := args[1].(*configuration.SystemConf)

	buyCoins := []string{}
	sortCoins := map[string]float32{}
	conditions := map[string]bool{}

	// fetch all weight coins for judging btc
	weights, err := utils.ReadLines(sysConf.WeightCsv)
	if err != nil {
		logrus.Error("read weight csv failed, err is ", err)
	}

	coints, err := database.GetAllCoint(db)
	if err != nil || len(coints) == 0 {
		logrus.Error("get cointegration from database error: ", err)
	}

	for _, coint := range coints {
		if containsBtc := strings.Contains(coint.Pair, utils.IndexCoin); containsBtc {
			sortCoins[coint.Pair] = utils.StringToFloat32(coint.Pvalue)
		}
	}

	cSorts := (&cSort{}).sortCoints(sortCoins)
	for i, weight := range weights {
		cointFlag := false
		for _, cSort := range cSorts {
			for _, coin := range strings.Split(cSort.Pair, "-") {
				if coin == utils.IndexCoin {
					if strings.Contains(cSort.Pair, weight) {
						cointFlag = true
					}
				}
			}
		}

		if !cointFlag {
			//TODO(wangxiaoyu), here has a bug.
			weights = append(weights[:i], weights[i+1:]...)
			logrus.Error("There is no correlation with BTC, the coin is ", weight)
		}
	}

	// for rising market
	btcRisingCondition := conditionUpMonitor(utils.IndexCoin, 1.0)

	sports := (&interfaces.MarketArgs{
		CurrencyPair: utils.IndexCoin,
		Interval:     utils.Now,
		Level:        utils.Level4Hour,
	}).SpotMarket()
	currentPrice := utils.StringToFloat64(sports[0][2])

	// fetch the 21 interval average of 4h
	averageArgs := index.Average{
		CurrencyPair: utils.IndexCoin,
		Level:        utils.Level4Hour,
		MA:           utils.MA21,
	}
	average21Per4h := averageArgs.Average(false)
	priceRisingCondition := currentPrice > average21Per4h

	for _, weight := range weights {
		// judgment depends on price average data
		conditions[weight] = conditionUpMonitor(weight, 1.0)
	}

	countUp := 0
	allUp := 0
	for _, condition := range conditions {
		if condition {
			countUp++
		}
		allUp++
	}

	// for falling market
	btcFallingCondition := conditionDownMonitor(utils.IndexCoin, 1.0)

	sports = (&interfaces.MarketArgs{
		CurrencyPair: utils.IndexCoin,
		Interval:     utils.Now,
		Level:        utils.Level4Hour,
	}).SpotMarket()
	currentPrice = utils.StringToFloat64(sports[0][2])

	// fetch the 21 interval average of 4h
	averageArgs = index.Average{
		CurrencyPair: utils.IndexCoin,
		Level:        utils.Level4Hour,
		MA:           utils.MA21,
	}
	average21Per4h = averageArgs.Average(false)
	priceFallingCondition := currentPrice < average21Per4h

	for _, weight := range weights {
		// judgment depends on price average data
		conditions[weight] = conditionDownMonitor(weight, 1.0)
	}

	countDown := 0
	allDown := 0
	for _, condition := range conditions {
		if condition {
			countDown++
		}
		allDown++
	}

	if float32(countUp)/float32(allUp) > 0.7 && btcRisingCondition && priceRisingCondition && averageDiff(utils.IndexCoin, utils.Level4Hour) {
		if tradeJugde(utils.IndexCoin, db) {
			buyCoins = append(buyCoins, utils.IndexCoin)
		}
	}

	if float32(countDown)/float32(allDown) > 0.7 && btcFallingCondition && priceFallingCondition && averageDiff(utils.IndexCoin, utils.Level4Hour) {
		if tradeJugde(utils.IndexCoin, db) {
			buyCoins = append(buyCoins, utils.IndexCoin)
		}
	}

	return buyCoins
}

func conditionUpMonitor(coin string, tenAverageDiff float64) bool {
	averageArgs := index.Average{
		CurrencyPair: coin,
		Level:        utils.Level4Hour,
		MA:           utils.MA21,
	}
	MA21Average := averageArgs.Average(false) > averageArgs.Average(true)*tenAverageDiff //4h的Average是增长的

	averageArgs.MA = utils.MA10
	MA10Average := averageArgs.Average(false) > averageArgs.Average(true)

	averageArgs.MA = utils.MA5
	MA15Average := averageArgs.Average(false) > averageArgs.Average(true)

	return MA21Average && MA10Average && MA15Average
}

func conditionDownMonitor(coin string, averageDiff float64) bool {
	averageArgs := index.Average{
		CurrencyPair: coin,
		Level:        utils.Level4Hour,
		MA:           utils.MA21,
	}
	MA21Average := averageArgs.Average(false)*averageDiff < averageArgs.Average(true) //4h的Average是增长的

	averageArgs.MA = utils.MA10
	MA10Average := averageArgs.Average(false) < averageArgs.Average(true)

	averageArgs.MA = utils.MA5
	MA15Average := averageArgs.Average(false) < averageArgs.Average(true)

	return MA21Average && MA10Average && MA15Average
}

// to judge this coin can be traded or not
func tradeJugde(coin string, db *gorm.DB) bool {
	inOrder := database.InOrder{
		Contract:  coin,
		Direction: "up",
	}

	record, err := inOrder.FetchOneInOrder(db)
	if err != nil && record == nil {
		return true
	} else {
		return false
	}

}

type cSort struct {
	Pair  string
	Value float32
}

func (*cSort) sortCoints(coints map[string]float32) []cSort {
	var cSorts []cSort
	for k, v := range coints {
		cSorts = append(cSorts, cSort{k, v})
	}

	sort.Slice(cSorts, func(i, j int) bool {
		return cSorts[i].Value < cSorts[j].Value // 升序
	})

	return cSorts
}

func averageDiff(coin string, level string) bool {
	var maValues []float64
	var max float64
	var min float64

	mas := []int{
		utils.MA5,
		utils.MA10,
		utils.MA21,
	}

	for _, ma := range mas {
		averageArgs := index.Average{
			CurrencyPair: coin,
			Level:        level,
			MA:           ma,
		}
		maValues = append(maValues, averageArgs.Average(false))
	}

	if len(maValues) != 3 {
		return false
	}

	for _, value := range maValues {
		if max == 0 {
			max = value
			min = value
		}

		if value >= max {
			max = value
		}

		if value < min {
			min = value
		}
	}

	return (max-min)/min > 0.03
}
