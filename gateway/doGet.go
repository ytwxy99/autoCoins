package gateway

import (
	"fmt"
	"net/http"

	"github.com/gateio/gateapi-go/v6"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	c "github.com/ytwxy99/autoCoins/client"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/utils"
)

func ReadLog(context *gin.Context, filePath string) {
	var content string
	logContent, _ := utils.ReadLines(filePath)
	if len(logContent) > 100 {
		for _, v := range logContent[len(logContent)-100 : len(logContent)] {
			content = content + "\n" + v
		}

	} else {
		for _, v := range logContent {
			content = content + "\n" + v
		}
	}

	context.String(http.StatusOK, content)
}

func ReadSold(context *gin.Context, db *gorm.DB) {
	var sum float32
	solds, err := database.GetAllSold(db)
	if err != nil {
		logrus.Error("get all solds err: ", err)
	}

	for _, sold := range solds {
		content := fmt.Sprintf("Sold detail: %s -> %s", sold.Contract, sold.Profit)
		fmt.Println(content)
		sum = sum + utils.StringToFloat32(sold.Relative_profit)
	}

	sums := fmt.Sprintf("Total sold profits is %s", sum)
	context.String(http.StatusOK, sums)
}

func ReadOrder(client *gateapi.APIClient, context *gin.Context, db *gorm.DB) {
	var contents string
	orders, err := database.GetAllOrder(db)
	if err != nil {
		logrus.Error("get all orders err: ", err)
	}

	for _, order := range orders {
		currentCoin, err := c.GetCurrencyPair(client, order.Contract)
		if err != nil {
			context.String(http.StatusInternalServerError, "Get last price failed: ", err)
		}
		lastPrice := utils.StringToFloat32(currentCoin[0].Last)
		priceDiff := (lastPrice - utils.StringToFloat32(order.Price)) / utils.StringToFloat32(order.Price) * 100
		contents = contents + fmt.Sprintf("order detail: coin -> %s, price -> %s, time -> %s, priceDiff -> %s \n", order.Contract, order.Price, order.CreatedAt, priceDiff)
	}

	context.String(http.StatusOK, contents)
}
