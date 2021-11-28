package db

import (
	"time"
)

type TxType uint8

const (
	UnknownTxType TxType = iota
	TokenTxType
	ContractTxType
)

type Tx struct {
	Hash         string    `gorm:"primaryKey"`
	BlockNumber  uint64    `gorm:"index:idx_tx_blk_num,sort:desc"`
	Timestamp    time.Time `gorm:"index:idx_tx_blk_timestamp,sort:desc"`
	From         string
	To           string
	Value        uint64
	EthTxType    uint8
	TxType       TxType
	GasUsed      uint64
	GasLimit     uint64
	InputDataHex string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
