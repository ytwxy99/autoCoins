package policy

import (
	"gorm.io/gorm"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/configuration"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
	"github.com/ytwxy99/autoCoins/utils/index"
)

type Umbrella struct{}

func (*Umbrella) Target(args ...interface{}) interface{} {
	db := args[0].(*gorm.DB)
	sysConf := args[1].(*configuration.SystemConf)

	sortCoins := map[string]float32{}
	buyCoins := []string{}
	conditions := map[string]bool{}

	//fetch all weight coins for judging bnb
	weights, err := utils.ReadLines(sysConf.Platform)
	if err != nil {
		logrus.Error("read platform csv failed, err is ", err)
	}

	coints, err := database.GetAllCoint(db)
	if err != nil || len(coints) == 0 {
		logrus.Error("get cointegration from database error: ", err)
	}

	for _, coint := range coints {
		if containsBtc := strings.Contains(coint.Pair, utils.IndexPlatformCoin); containsBtc {
			sortCoins[coint.Pair] = utils.StringToFloat32(coint.Pvalue)
		}
	}

	cSorts := (&cSort{}).sortCoints(sortCoins)
	for i, weight := range weights {
		cointFlag := false
		for _, cSort := range cSorts {
			for _, coin := range strings.Split(cSort.Pair, "-") {
				if coin == utils.IndexPlatformCoin {
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
	btcRisingCondition := conditionUpMonitor(utils.IndexPlatformCoin, 1.0)

	sports := (&interfaces.MarketArgs{
		CurrencyPair: utils.IndexPlatformCoin,
		Interval:     utils.Now,
		Level:        utils.Level4Hour,
	}).SpotMarket()
	currentPrice := utils.StringToFloat64(sports[0][2])

	// fetch the 21 interval average of 4h
	averageArgs := index.Average{
		CurrencyPair: utils.IndexPlatformCoin,
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
	btcFallingCondition := conditionDownMonitor(utils.IndexPlatformCoin, 1.0)

	sports = (&interfaces.MarketArgs{
		CurrencyPair: utils.IndexPlatformCoin,
		Interval:     utils.Now,
		Level:        utils.Level4Hour,
	}).SpotMarket()
	currentPrice = utils.StringToFloat64(sports[0][2])

	// fetch the 21 interval average of 4h
	averageArgs = index.Average{
		CurrencyPair: utils.IndexPlatformCoin,
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

	if float32(countUp)/float32(allUp) > 0.95 && btcRisingCondition && priceRisingCondition && averageDiff(utils.IndexPlatformCoin, utils.Level4Hour) {
		if tradeJugde(utils.IndexPlatformCoin, db) {
			buyCoins = append(buyCoins, utils.IndexPlatformCoin)
		}
	}

	if float32(countDown)/float32(allDown) > 0.95 && btcFallingCondition && priceFallingCondition && averageDiff(utils.IndexPlatformCoin, utils.Level4Hour) {
		if tradeJugde(utils.IndexPlatformCoin, db) {
			buyCoins = append(buyCoins, utils.IndexPlatformCoin)
		}
	}

	return buyCoins

	return nil
}
