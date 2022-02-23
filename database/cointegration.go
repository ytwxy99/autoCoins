package database

import (
	"gorm.io/gorm"
)

// get all cointegrations or spicified limit
func GetAllCoint(db *gorm.DB) ([]Cointegration, error) {
	var coints []Cointegration
	tx := db.Find(&coints)

	return coints, tx.Error
}
