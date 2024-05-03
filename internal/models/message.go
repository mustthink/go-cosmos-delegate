package models

type DelegateMessage struct {
	ID               uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	TransactionID    uint64 `gorm:"column:transaction_id;type:bigint;not null"`
	DelegatorAddress string `gorm:"column:delegator_address;type:varchar(64);not null"`
	ValidatorAddress string `gorm:"column:validator_address;type:varchar(64);not null"`
	Amount           int64  `gorm:"column:amount;type:bigint;not null"`
	Currency         string `gorm:"column:currency;type:varchar(64);not null"`
}
