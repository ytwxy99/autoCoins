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

// find buy point by doing cointegration
func (*Cointegration) Target(args ...interface{}) interface{} {
	statistics := make(map[string]float64)
	//convert specified type
	client := args[0].(*gateapi.APIClient)
	db := args[1].(*gorm.DB)

	coints, err := database.GetAllCoint(db)
	if err != nil {
		logrus.Error("get cointegration from database error:", err)
	}

	coints = removeDuplicate(coints)
	logrus.Info("P-Value less than 0.000001, totally: ", len(coints))
	for _, value := range coints {
		pairs := strings.Split(value.Pair, "-")
		k0 := interfaces.K(client, pairs[0], -3, "1d")
		k1 := interfaces.K(client, pairs[1], -3, "1d")
		price0 := (utils.StringToFloat32(k0[3][2]) - utils.StringToFloat32(k0[2][2])) / utils.StringToFloat32(k0[2][2])
		price1 := (utils.StringToFloat32(k1[3][2]) - utils.StringToFloat32(k1[2][2])) / utils.StringToFloat32(k1[2][2])

		if math.Abs(float64(price0-price1)) >= 0.2 && (price0 > 0 || price1 > 0) {
			if _, ok := statistics[value.Pair]; !ok {
				statistics[value.Pair] = math.Abs(float64(price0 - price1))
			}
		}
	}

	if len(statistics) != 0 {
		for k, v := range statistics {
			logrus.Info("Find cointegration buy point:", k, v)
		}
	}

	return nil
}
