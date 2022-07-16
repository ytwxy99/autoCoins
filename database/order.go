package database

import (
	"context"
	"errors"

	"github.com/ytwxy99/autocoins/pkg/utils"
)

// AddOrder add one order
func (order *Order) AddOrder(ctx context.Context) error {
	tx := utils.GetDBContext(ctx).Create(order)
	return tx.Error
}

// FetchOneOrder fetch one order by contract and ditection
func (order Order) FetchOneOrder(ctx context.Context) (*Order, error) {
	utils.GetDBContext(ctx).Table("orders").
		Where("contract = ? AND direction = ?", order.Contract, order.Direction).First(&order)

	if order.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &order, nil
}

// UpdateOrder update order
func (order *Order) UpdateOrder(ctx context.Context) error {
	tx := utils.GetDBContext(ctx).Model(&Order{}).Where("price > ?", 10).Updates(order)
	return tx.Error
}

// DeleteOrder delete order
func (order *Order) DeleteOrder(ctx context.Context) error {
	tx := utils.GetDBContext(ctx).Delete(order)
	return tx.Error
}

// GetAllOrder get all order
func GetAllOrder(ctx context.Context) ([]Order, error) {
	var orders []Order
	tx := utils.GetDBContext(ctx).Find(&orders)
	return orders, tx.Error
}
