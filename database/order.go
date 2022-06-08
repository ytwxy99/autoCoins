package database

import (
	"errors"

	"gorm.io/gorm"
)

// add one order
func (order *Order) AddOrder(db *gorm.DB) error {
	tx := db.Create(order)
	return tx.Error
}

// fetch one order by contract and ditection
func (order Order) FetchOneOrder(db *gorm.DB) (*Order, error) {
	db.Table("orders").
		Where("contract = ? AND direction = ?", order.Contract, order.Direction).First(&order)

	if order.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &order, nil
}

// update order
func (order *Order) UpdateOrder(db *gorm.DB) error {
	tx := db.Model(&Order{}).Where("price > ?", 10).Updates(order)
	return tx.Error
}

// delete order
func (order *Order) DeleteOrder(db *gorm.DB) error {
	tx := db.Delete(order)
	return tx.Error
}

// get all order
func GetAllOrder(db *gorm.DB) ([]Order, error) {
	var orders []Order
	tx := db.Find(&orders)
	return orders, tx.Error
}
