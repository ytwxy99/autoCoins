package database

import (
	"gorm.io/gorm"
)

// add one history_day
func (historyDay *HistoryDay) AddHistoryDay(db *gorm.DB) error {
	tx := db.Create(historyDay)
	return tx.Error
}

// get all history_day
func GetAllHistoryDay(db *gorm.DB) ([]HistoryDay, error) {
	var historyDays []HistoryDay
	tx := db.Find(&historyDays)
	return historyDays, tx.Error
}