package db

import (
	"time"
)

type Block struct {
	Number     uint64    `gorm:"primaryKey,sort:desc"`
	Hash       string    `gorm:"index:idx_blk_hash"`
	ParentHash string    `gorm:"index:idx_blk_hash"`
	Timestamp  time.Time `gorm:"index:blk_timestamp,sort:desc"`
	GasUsed    uint64
	GasLimit   uint64
	Raw        string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
