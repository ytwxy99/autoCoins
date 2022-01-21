package database

import (
	"errors"
	"gorm.io/gorm"
)

// add one sold
func (sold *Sold) AddSold(db *gorm.DB) error {
	tx := db.Create(sold)
	return tx.Error
}

// fetch one sold by contract and ditection
func (sold Sold) FetchOneSold(db *gorm.DB) (*Sold, error) {
	db.Table("sold").
		Where("contract = ?", sold.Contract).First(&sold)

	if sold.ID == 0 {
		return nil, errors.New("record not found")
	}
	return &sold, nil
}

// get all sold
func GetAllSold(db *gorm.DB) ([]Sold, error) {
	var solds []Sold
	tx := db.Find(&solds)
	return solds, tx.Error
}
