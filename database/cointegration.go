package database

import (
	"context"

	"github.com/ytwxy99/autocoins/pkg/utils"
)

// GetAllCoint get all cointegrations or spicified limit
func GetAllCoint(ctx context.Context) ([]Cointegration, error) {
	var coints []Cointegration
	tx := utils.GetDBContext(ctx).Find(&coints)

	return coints, tx.Error
}
