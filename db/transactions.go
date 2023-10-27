package db

import (
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
	CreateTransactions(transactions *models.Transaction) (*models.Transaction, error)
}

func (t *transactions) CreateTransactions(transactions *models.Transaction) (*models.Transaction, error) {
	err := t.DB.Create(transactions).Error
	if err != nil {
		return nil, fmt.Errorf("could not create transactions: %v", err)
	}
	return transactions, nil
}
