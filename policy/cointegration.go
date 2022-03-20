package policy

import (
	"gorm.io/gorm"
	"math"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
	"github.com/ytwxy99/autoCoins/utils/index"
)

type Cointegration struct{}

// duplicate remove the result of cointegration
func removeDuplicate(coints []database.Cointegration) []database.Cointegration {
	var recordCoints []string
	var newCoints []database.Cointegration
	for _, coint := range coints {
		isExist := false
		pValue := utils.StringToFloat32(coint.Pvalue)
		sPair := strings.Split(coint.Pair, "-")
		compare := sPair[1] + "-" + sPair[0]

		for _, recordCoint := range recordCoints {
			if recordCoint == coint.Pair || recordCoint == compare {
				isExist = true
			}
		}

		if !isExist {
			recordCoints = append(recordCoints, coint.Pair)
			if pValue <= 0.000001 {
				newCoints = append(newCoints, coint)
			}
		}

	}

	return newCoints
}

// MACD condition judge
func macdJudge(coin string, interval int, level string) bool {
	k4hValues := interfaces.Market(coin, interval, level)
	if k4hValues != nil {
		k4hMacds := index.GetMacd(k4hValues, 12, 26, 9)
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

// can be traded or not
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

// find buy point by doing cointegration
func (*Cointegration) Target(args ...interface{}) interface{} {
	buyCoins := []string{}
	priceDiff := make(map[string]float32)
	duplicates := make(map[string]int32)
	statistics := make(map[string]float64)

	//convert specified type
	db := args[0].(*gorm.DB)

	coints, err := database.GetAllCoint(db)
	if err != nil {
		logrus.Error("get cointegration from database error: %v", err)
	}

	//TODO(wangxiaoyu), maybe have a bug here.
	coints = removeDuplicate(coints)

	for _, value := range coints {
		pairs := strings.Split(value.Pair, "-")
		if len(pairs) != 2 {
			logrus.Error("coints pair record error")
			continue
		}

		k0 := interfaces.Market(pairs[0], -3, "1d")
		k1 := interfaces.Market(pairs[1], -3, "1d")
		if k0 == nil || k1 == nil {
			// get k data failed
			continue
		}

		diff0 := utils.PriceDiffPercent(k0[3][2], k0[2][2])
		diff1 := utils.PriceDiffPercent(k1[3][2], k1[2][2])
		priceDiff[pairs[0]] = diff0
		priceDiff[pairs[1]] = diff1

		if math.Abs(float64(diff0-diff1)) >= 0.2 && (diff0 > 0 || diff1 > 0) {
			if _, ok := statistics[value.Pair]; !ok {
				statistics[value.Pair] = math.Abs(float64(diff0 - diff1))
			}
			contractNames := strings.Split(value.Pair, "-")
			for _, name := range contractNames {
				if _, ok := duplicates[name]; ok {
					duplicates[name]++
				} else {
					duplicates[name] = 1
				}
			}
		}
	}

	if len(statistics) != 0 {
		for k, v := range statistics {
			pairs := strings.Split(k, "-")
			if len(pairs) == 2 && (duplicates[pairs[0]] > 1 || duplicates[pairs[1]] > 1) {
				logrus.Warnf("Suspected to be buying point: %v", k, " price diff: %v", v)
				continue
			}

			//TODO(wangxiaoyu), need to optimize buying points.
			paris := strings.Split(k, "-")
			if priceDiff[paris[0]] > priceDiff[paris[1]] {
				if macdJudge(paris[1], -100, "4h") && macdJudge(paris[1], -10, "15m") {
					if tradeJugde(paris[1], db) {
						tradeDetail := database.TradeDetail{
							Contract:  paris[1],
							CointPair: paris[0],
						}
						err = (&tradeDetail).AddTradeDetail(db)
						if err != nil {
							logrus.Errorf("add TradeDetail error in buy point : %v , inOrder is %v:", err, tradeDetail)
							continue
						}

						buyCoins = append(buyCoins, paris[1])
						logrus.Info("Find cointegration buy point: %v", paris[1], " contract pairs: %v", k, " price diff: %v", priceDiff[paris[1]])
					}
				}
			} else {
				if macdJudge(paris[0], -100, "4h") && macdJudge(paris[1], -10, "15m") {
					if tradeJugde(paris[0], db) {
						tradeDetail := database.TradeDetail{
							Contract:  paris[0],
							CointPair: paris[1],
						}
						err = (&tradeDetail).AddTradeDetail(db)
						if err != nil {
							logrus.Errorf("add TradeDetail error in buy point : %v , inOrder is : %v", err, tradeDetail)
							continue
						}

						buyCoins = append(buyCoins, paris[0])
						logrus.Info("Find cointegration buy point: %v", paris[0], " contract pairs: %v", k, " price diff: %v", priceDiff[paris[0]])
					}
				}
			}
		}

		return buyCoins
	} else {
		return buyCoins
	}
}
