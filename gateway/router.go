package gateway

import (
	"gorm.io/gorm"

	"github.com/gateio/gateapi-go/v6"
	"github.com/gin-gonic/gin"

	"github.com/ytwxy99/autoCoins/configuration"
)

func Router(client *gateapi.APIClient, engine *gin.Engine, sysConf *configuration.SystemConf, db *gorm.DB) {
	engine.GET("/", func(context *gin.Context) {
		ReadLog(context, sysConf.LogPath)
	})
	engine.GET("/solds", func(context *gin.Context) {
		ReadSold(context, db)
	})
	engine.GET("/orders", func(context *gin.Context) {
		ReadOrder(client, context, db)
	})

	// run gateway service
	engine.Run(":80")
}
