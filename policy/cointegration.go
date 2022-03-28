package policy

import (
	"gorm.io/gorm"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
	"github.com/ytwxy99/autoCoins/utils/index"
)

type Cointegration struct{}

// cointegration policy
func (*Cointegration) Target(args ...interface{}) interface{} {
	db := args[0].(*gorm.DB)
	buyCoins := []string{}
	sortCoins := map[string]float32{}

	btcMarket := (&interfaces.MarketArgs{
		CurrencyPair: utils.IndexCoin,
		Interval:     -1,
		Level:        utils.Level4Hour,
	}).Market()

	diffBtc1 := utils.Compare(btcMarket[len(btcMarket)-1][2], btcMarket[len(btcMarket)-2][2], 0, 1.05)
	diffBtc2 := utils.Compare(btcMarket[len(btcMarket)-3][2], btcMarket[len(btcMarket)-3][2], 0, 1.01)

	if diffBtc1 && !diffBtc2 {
		// 当前4个小时涨幅5%，且上个4个小时涨幅小于1%
		coints, err := database.GetAllCoint(db)
		if err != nil || len(coints) == 0 {
			logrus.Error("get cointegration from database error: %v", err)
		}

		for _, coint := range coints {
			if containsBtc := strings.Contains(coint.Pair, utils.IndexCoin); containsBtc {
				sortCoins[coint.Pair] = utils.StringToFloat32(coint.Pvalue)
			}
		}

		cSorts := (&cSort{}).sortCoints(sortCoins)
		for _, cSort := range cSorts {
			pairs := strings.Split(cSort.Pair, "-")
			for _, p := range pairs {
				if utils.IndexCoin == p {
					continue
				}

				btcMarket := (&interfaces.MarketArgs{
					CurrencyPair: p,
					Interval:     -1,
					Level:        utils.Level4Hour,
				}).Market()

				diffSubCoin1 := utils.Compare(btcMarket[len(btcMarket)-1][2], btcMarket[len(btcMarket)-2][2], 0, 1.05)
				diffSubCoin2 := utils.Compare(btcMarket[len(btcMarket)-3][2], btcMarket[len(btcMarket)-3][2], 0, 1.05)

				if !diffSubCoin1 && !diffSubCoin2 {
					if tradeJugde(p, db) {
						buyCoins = append(buyCoins, p)
						logrus.Info("Find cointegration policy desired coin: %v", p)
					}
				}
			}
		}
	}

	return buyCoins
}

func macdJudge(marketArgs *interfaces.MarketArgs) bool {
	k4hValues := marketArgs.Market()
	if k4hValues != nil {
		macdArgs := index.DefaultMacdArgs()
		k4hMacds := macdArgs.GetMacd(k4hValues)
		nowK4h := len(k4hMacds) - 1
		if nowK4h < 5 {
			return false
		}

		macdValue := utils.StringToFloat32(k4hMacds[nowK4h]["macd"])
		if macdValue > 0 {
			return true
		}
	}

	return false
}

// to judge this coin can be trade or not
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
