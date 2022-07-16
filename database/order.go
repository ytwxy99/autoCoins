package database

import (
	"errors"

	"gorm.io/gorm"
)

// AddOrder add one order
func (order *Order) AddOrder(db *gorm.DB) error {
	tx := db.Create(order)
	return tx.Error
}

// FetchOneOrder fetch one order by contract and ditection
func (order Order) FetchOneOrder(db *gorm.DB) (*Order, error) {
	db.Table("orders").
		Where("contract = ? AND direction = ?", order.Contract, order.Direction).First(&order)

	if order.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &order, nil
}

// UpdateOrder update order
func (order *Order) UpdateOrder(db *gorm.DB) error {
	tx := db.Model(&Order{}).Where("price > ?", 10).Updates(order)
	return tx.Error
}

// DeleteOrder delete order
func (order *Order) DeleteOrder(db *gorm.DB) error {
	tx := db.Delete(order)
	return tx.Error
}

// GetAllOrder get all order
func GetAllOrder(db *gorm.DB) ([]Order, error) {
	var orders []Order
	tx := db.Find(&orders)
	return orders, tx.Error
}
