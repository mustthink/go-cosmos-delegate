package db

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mustthink/go-cosmos-delegate/internal/config"
	"github.com/mustthink/go-cosmos-delegate/internal/models"
)

// Storage is a storage layer for the news service.
type (
	database struct {
		db *gorm.DB
	}

	Storage interface {
		CreateTransaction(ctx context.Context, transaction models.Transaction) (uint64, error)
		CreateDelegateMessages(ctx context.Context, transactionID uint64, messages []models.DelegateMessage) error
		Close() error
	}
)

// New creates a new storage layer.
func New(cfg config.Database) (Storage, error) {
	const op = "storage.New"

	dialector := postgres.Open(cfg.Uri())
	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &database{db: db}, nil
}

// MustNew creates a new storage layer and panics if an error occurs.
func MustNew(cfg config.Database) Storage {
	st, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return st
}

// Close closes the storage layer.
func (s *database) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
