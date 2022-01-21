package database

import (
	"errors"
	"gorm.io/gorm"
)

// add one order
func (inOrder *InOrder) AddInOrder(db *gorm.DB) error {
	tx := db.Create(inOrder)
	return tx.Error
}

// fetch one order by contract and ditection
func (inOrder InOrder) FetchOneInOrder(db *gorm.DB) (*InOrder, error) {
	db.Table("inOrder").
		Where("contract = ? AND direction = ?", inOrder.Contract, inOrder.Direction).First(&inOrder)

	if inOrder.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &inOrder, nil
}

// update order
func (inOrder *InOrder) UpdateInOrder(db *gorm.DB) error {
	tx := db.Model(&InOrder{}).Where("price > ?", 10).Updates(inOrder)
	return tx.Error
}

// delete order
func (inOrder *InOrder) DeleteInOrder(db *gorm.DB) error {
	tx := db.Delete(inOrder)
	return tx.Error
}
