package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/ytwxy99/autoCoins/configuration"
	"gorm.io/gorm"
)

func Router(engine *gin.Engine, sysConf *configuration.SystemConf, db *gorm.DB) {
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
