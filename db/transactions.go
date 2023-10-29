package db

import (
	"errors"
	"escrolla-api/models"
	"fmt"
	"gorm.io/gorm"
)

type transactions struct {
	DB *gorm.DB
}

func NewTransactions(db *GormDB) TransactionsRepo {
	return &transactions{db.DB}
}

type TransactionsRepo interface {
	CreateOrder(order *models.Order) (*models.Order, error)
	UpdateOrderStatus(reference string, newStatus string) error
}

func (t *transactions) CreateOrder(order *models.Order) (*models.Order, error) {
	err := t.DB.Create(order).Error
	if err != nil {
		return nil, fmt.Errorf("could not create order: %v", err)
	}
	return order, nil
}

// UpdateOrderStatus updates the order status based on the reference.
func (t *transactions) UpdateOrderStatus(reference string, newStatus string) error {
	var order models.Order

	// Find the order by its reference
	if err := t.DB.Where("id = ?", reference).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%v order not found")
		}
		return err
	}

	// Update the order status
	order.PaymentStatus = newStatus

	// Save the updated order to the database
	if err := t.DB.Save(&order).Error; err != nil {
		return err
	}

	return nil
}
