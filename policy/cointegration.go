package policy

import (
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ytwxy99/autoCoins/configuration"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/utils"
	"github.com/ytwxy99/autoCoins/utils/index"
)

type Cointegration struct{}

// cointegration policy
func (*Cointegration) Target(args ...interface{}) interface{} {
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
			weights = append(weights[:i], weights[i+1:]...)
			logrus.Error("There is no correlation with BTC, the coin is ", weight)
		}
	}

	// monitor btc
	btcCondition := conditionMonitor(utils.IndexCoin, 1.01)

	for _, weight := range weights {
		// judgment depends on price average data
		conditions[weight] = conditionMonitor(weight, 1.001)
	}

	count := 0
	all := 0
	for _, condition := range conditions {
		if condition {
			count++
		}
		all++
	}

	if float32(count)/float32(all) > 0.7 && btcCondition {
		if tradeJugde(utils.IndexCoin, db) {
			buyCoins = append(buyCoins, utils.IndexCoin)
		}
	}

	return buyCoins
}

func conditionMonitor(coin string, tenAverageDiff float64) bool {
	averageArgs := index.Average{
		CurrencyPair: coin,
		Intervel:     utils.Thirty,
		Level:        utils.Level4Hour,
	}
	btcThirtyAverage := averageArgs.Average(false) >= averageArgs.Average(true)*tenAverageDiff //4h的Average是增长的

	averageArgs.Intervel = utils.Ten
	btcTenAverage := averageArgs.Average(false) > averageArgs.Average(true)

	averageArgs.Intervel = utils.Five
	btcFiveAverage := averageArgs.Average(false) > averageArgs.Average(true)

	averageArgs.Intervel = utils.Thirty
	btcConditionA := averageArgs.Average(false) <= averageArgs.Average(true)*1.001

	averageArgs.Intervel = utils.Ten
	btcConditionB := averageArgs.Average(false) > 0

	averageArgs.Intervel = utils.Five
	btcConditionC := averageArgs.Average(false) > 0

	return btcThirtyAverage && btcTenAverage && btcFiveAverage && btcConditionA && btcConditionB && btcConditionC
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
