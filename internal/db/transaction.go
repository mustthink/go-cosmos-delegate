package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/mustthink/go-cosmos-delegate/internal/models"
)

func (s *database) CreateTransaction(ctx context.Context, transaction models.Transaction) (uint64, error) {
	err := s.db.WithContext(ctx).
		Model(models.Transaction{}).
		Where("external_id = ?", transaction.ExternalID).
		Find(&transaction).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) || transaction.ID != 0 {
		return 0, err
	}

	err = s.db.WithContext(ctx).
		Model(models.Transaction{}).
		Create(&transaction).
		Error
	return transaction.ID, err
}

func (s *database) CreateDelegateMessages(ctx context.Context, transactionID uint64, messages []models.DelegateMessage) error {
	for i := range messages {
		messages[i].TransactionID = transactionID
	}

	return s.db.WithContext(ctx).
		Model(models.DelegateMessage{}).
		Create(messages).
		Error
}
