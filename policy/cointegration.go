package policy

import (
	"gorm.io/gorm"
	"math"
	"strings"

	"github.com/gateio/gateapi-go/v6"
	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/interfaces"
	"github.com/ytwxy99/autoCoins/utils"
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
func macdJudge(client *gateapi.APIClient, coin string, level string) bool {
	// judge depends '4 hours' data
	k4hValues := interfaces.K(client, coin, -100, level)
	if k4hValues != nil {
		k4hMacds := GetMacd(k4hValues, 12, 26, 9)
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
	client := args[0].(*gateapi.APIClient)
	db := args[1].(*gorm.DB)

	coints, err := database.GetAllCoint(db)
	if err != nil {
		logrus.Error("get cointegration from database error:", err)
	}

	coints = removeDuplicate(coints)
	//logrus.Info("P-Value less than 0.000001, totally: ", len(coints))
	for _, value := range coints {
		pairs := strings.Split(value.Pair, "-")
		k0 := interfaces.K(client, pairs[0], -3, "1d")
		k1 := interfaces.K(client, pairs[1], -3, "1d")
		price0 := (utils.StringToFloat32(k0[3][2]) - utils.StringToFloat32(k0[2][2])) / utils.StringToFloat32(k0[2][2])
		price1 := (utils.StringToFloat32(k1[3][2]) - utils.StringToFloat32(k1[2][2])) / utils.StringToFloat32(k1[2][2])
		priceDiff[pairs[0]] = price0
		priceDiff[pairs[0]] = price1

		if math.Abs(float64(price0-price1)) >= 0.2 && (price0 > 0 || price1 > 0) {
			if _, ok := statistics[value.Pair]; !ok {
				statistics[value.Pair] = math.Abs(float64(price0 - price1))
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
			for _, name := range strings.Split(k, "-") {
				if duplicates[name] > 1 {
					logrus.Warn("Suspected to be buying point:", k, " price diff: ", v)
					continue
				}
			}

			//TODO(wangxiaoyu), need to optimize buying points.
			paris := strings.Split(k, "-")
			if priceDiff[paris[0]] > priceDiff[paris[1]] {
				if macdJudge(client, paris[1], "4h") && macdJudge(client, paris[1], "15m") {
					if tradeJugde(paris[1], db) {
						buyCoins = append(buyCoins, paris[1])
						logrus.Info("Find cointegration buy point:", paris[0], " contract pairs:", k, " price diff:", priceDiff[paris[1]])
					}
				}
			} else {
				if macdJudge(client, paris[0], "4h") && macdJudge(client, paris[1], "15m") {
					if tradeJugde(paris[0], db) {
						buyCoins = append(buyCoins, paris[0])
						logrus.Info("Find cointegration buy point:", paris[1], " contract pairs:", k, " price diff:", priceDiff[paris[0]])
					}
				}
			}
		}

		return buyCoins
	} else {
		return buyCoins
	}
}
