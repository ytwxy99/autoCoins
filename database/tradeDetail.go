package database

import (
	"errors"

	"gorm.io/gorm"
)

// add one TradeDetail
func (tradeDetail *TradeDetail) AddTradeDetail(db *gorm.DB) error {
	tx := db.Create(tradeDetail)
	return tx.Error
}

// fetch one TradeDetail
func (tradeDetail TradeDetail) FetchOneTradeDetail(db *gorm.DB) (*TradeDetail, error) {
	db.Table("inOrder").
		Where("contract = ? AND coint_pair = ?", tradeDetail.Contract, tradeDetail.CointPair).First(&tradeDetail)

	if tradeDetail.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &tradeDetail, nil
}

// delete TradeDetail
func (tradeDetail *TradeDetail) DeleteTradeDetail(db *gorm.DB) error {
	tx := db.Where("contract = ? and coint_pair = ?", tradeDetail.Contract, tradeDetail.CointPair).Delete(tradeDetail)
	return tx.Error
}
