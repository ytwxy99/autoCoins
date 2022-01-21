package gateway

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ytwxy99/autoCoins/database"
	"github.com/ytwxy99/autoCoins/utils"
	"gorm.io/gorm"
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

func ReadOrder(context *gin.Context, db *gorm.DB) {
	var contents string
	orders, err := database.GetAllOrder(db)
	if err != nil {
		logrus.Error("get all orders err: ", err)
	}

	for _, order := range orders {
		content := fmt.Sprintf("order detail: coin -> %s, price -> %s, time -> %s", order.Contract, order.Price, order.CreatedAt)
		fmt.Println(content)
		contents = contents + "\n" + content
	}

	context.String(http.StatusOK, contents)
}
