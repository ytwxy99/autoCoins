package utils

import (
	"context"
	"gorm.io/gorm"

	"github.com/ytwxy99/autocoins/pkg/configuration"
)

var ctx context.Context

type SystemContext struct {
	AuthConf   *configuration.GateAPIV4
	SystemConf *configuration.SystemConf
	Database   *gorm.DB
}

//InitCtx initial context
func InitCtx() context.Context {
	ctx = context.Background()
	return ctx
}

//SetContextValue set k-v to context
func SetContextValue(ctx context.Context, key interface{}, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetSystemConfContext(ctx context.Context) *configuration.SystemConf {
	return ctx.Value("ctxMetadata").(SystemContext).SystemConf
}

func GetDBContext(ctx context.Context) *gorm.DB {
	return ctx.Value("ctxMetadata").(SystemContext).Database
}
