package gateway

import (
	"context"
	"github.com/gin-gonic/gin"
)

func Router(ctx context.Context, engine *gin.Engine) {

	engine.GET("/", func(ginCtx *gin.Context) {
		ReadLog(ctx, ginCtx)
	})
	engine.GET("/solds", func(ginCtx *gin.Context) {
		ReadSold(ctx, ginCtx)
	})
	engine.GET("/orders", func(ginCtx *gin.Context) {
		ReadOrder(ctx, ginCtx)
	})

	// run gateway service
	engine.Run(":80")
}
