package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ytwxy99/autocoins/database"
	"github.com/ytwxy99/autocoins/pkg/client"
	"github.com/ytwxy99/autocoins/pkg/utils"
)

func ReadLog(ctx context.Context, ginCtx *gin.Context) {
	var content string
	logContent, _ := utils.ReadLines(utils.GetSystemConfContext(ctx).LogPath)
	if len(logContent) > 100 {
		for _, v := range logContent[len(logContent)-100 : len(logContent)] {
			content = content + "\n" + v
		}

	} else {
		for _, v := range logContent {
			content = content + "\n" + v
		}
	}

	ginCtx.String(http.StatusOK, content)
}

func ReadSold(ctx context.Context, ginCtx *gin.Context) {
	var sum float32
	solds, err := database.GetAllSold(ctx)
	if err != nil {
		logrus.Error("get all solds err: %v", err)
	}

	for _, sold := range solds {
		content := fmt.Sprintf("Sold detail: %s -> %s", sold.Contract, sold.Profit)
		fmt.Println(content)
		sum = sum + utils.StringToFloat32(sold.Relative_profit)
	}

	sums := fmt.Sprintf("Total sold profits is %s", sum)
	ginCtx.String(http.StatusOK, sums)
}

func ReadOrder(ctx context.Context, ginCtx *gin.Context) {
	var contents string
	orders, err := database.GetAllOrder(ctx)
	if err != nil {
		logrus.Error("get all orders err: %v", err)
	}

	for _, order := range orders {
		currentCoin, err := client.GetCurrencyPair(order.Contract)
		if err != nil {
			ginCtx.String(http.StatusInternalServerError, "Get last price failed: ", err)
		}

		priceDiff := utils.PriceDiffPercent(currentCoin[0].Last, order.Price)
		contents = contents + fmt.Sprintf("order detail: coin -> %s, price -> %s, time -> %s, priceDiff -> %s \n", order.Contract, order.Price, order.CreatedAt, priceDiff)
	}

	ginCtx.String(http.StatusOK, contents)
}
