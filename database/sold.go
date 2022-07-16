package database

import (
	"context"
	"errors"

	"github.com/ytwxy99/autocoins/pkg/utils"
)

// AddSold add one sold
func (sold *Sold) AddSold(ctx context.Context) error {
	tx := utils.GetDBContext(ctx).Create(sold)
	return tx.Error
}

// FetchOneSold fetch one sold by contract and ditection
func (sold Sold) FetchOneSold(ctx context.Context) (*Sold, error) {
	utils.GetDBContext(ctx).Table("sold").
		Where("contract = ?", sold.Contract).First(&sold)

	if sold.ID == 0 {
		return nil, errors.New("record not found")
	}
	return &sold, nil
}

// GetAllSold get all sold
func GetAllSold(ctx context.Context) ([]Sold, error) {
	var solds []Sold
	tx := utils.GetDBContext(ctx).Find(&solds)
	return solds, tx.Error
}
