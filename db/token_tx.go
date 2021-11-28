package db

import (
	"time"
)

type TokenTx struct {
	ID              uint      `gorm:"primaryKey"`
	TxHash          string    `gorm:"index:idx_ttx_hash"`
	BlockNumber     uint64    `gorm:"index:idx_ttx_blk_num,sort:desc"`
	Timestamp       time.Time `gorm:"index:idx_ttx_blk_timestamp,sort:desc"`
	From            string
	To              string
	Value           uint64
	Decimals        uint8
	ContractAddress string
	Symbol          string
	GasUsed         uint64
	GasLimit        uint64
	InputDataHex    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
