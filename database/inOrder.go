package database

import (
	"context"
	"errors"

	"github.com/ytwxy99/autocoins/pkg/utils"
)

// AddInOrder add one order
func (inOrder *InOrder) AddInOrder(ctx context.Context) error {
	tx := utils.GetDBContext(ctx).Create(inOrder)
	return tx.Error
}

// FetchOneInOrder fetch one order by contract and ditection
func (inOrder InOrder) FetchOneInOrder(ctx context.Context) (*InOrder, error) {
	utils.GetDBContext(ctx).Table("inOrder").
		Where("contract = ? AND direction = ?", inOrder.Contract, inOrder.Direction).First(&inOrder)

	if inOrder.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &inOrder, nil
}

// UpdateInOrder update order
func (inOrder *InOrder) UpdateInOrder(ctx context.Context) error {
	tx := utils.GetDBContext(ctx).Model(&InOrder{}).Where("price > ?", 10).Updates(inOrder)
	return tx.Error
}

// DeleteInOrder delete order
func (inOrder *InOrder) DeleteInOrder(ctx context.Context) error {
	tx := utils.GetDBContext(ctx).Where("contract = ? and direction = ?", inOrder.Contract, inOrder.Contract).Delete(inOrder)
	return tx.Error
}
