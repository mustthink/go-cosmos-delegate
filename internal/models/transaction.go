package models

import "time"

type Transaction struct {
	ID         uint64    `gorm:"type:bigserial; primaryKey; autoIncrement"`
	ExternalID string    `gorm:"uniqueIndex; type:varchar(64); not null"`
	BlockID    int64     `gorm:"type:bigint; not null"`
	Timestamp  time.Time `gorm:"type:timestamp; not null"`

	Messages []DelegateMessage `gorm:"-"`
}
