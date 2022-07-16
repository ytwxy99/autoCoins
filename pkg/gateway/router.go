package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ytwxy99/autocoins/pkg/utils"
)

func Router(ctx context.Context, engine *gin.Engine) {
	sysConf := utils.GetSystemConfContext(ctx)
	db := utils.GetDBContext(ctx)

	engine.GET("/", func(context *gin.Context) {
		ReadLog(context, sysConf.LogPath)
	})
	engine.GET("/solds", func(context *gin.Context) {
		ReadSold(context, db)
	})
	engine.GET("/orders", func(context *gin.Context) {
		ReadOrder(context, db)
	})

	// run gateway service
	engine.Run(":80")
}
