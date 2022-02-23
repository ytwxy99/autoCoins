package policy

import (
	"fmt"
	"gorm.io/gorm"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/utils"
)

type Cointegration struct{}

// find buy point by doing cointegration
func (*Cointegration) Target(args ...interface{}) interface{} {
	//convert specified type
	db := args[0].(*gorm.DB)

	coints, err := database.GetAllCoint(db)
	if err != nil {
		logrus.Error("get cointegration from database error:", err)
	}

	coints = removeDuplicate(coints)
	for _, value := range coints {
		fmt.Println(value.Pair, value.Pvalue)
	}
	logrus.Info("P-Value less than 0.000001, totally: ", len(coints))

	return nil
}

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
