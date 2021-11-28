package types

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type TxBlockDetails struct {
	Num  uint64      `json:"num"`
	Hash common.Hash `json:"hash"`
}

type TxToken struct {
	Symbol       string `json:"symbol"`
	TokenAddress string `json:"tokenAddress"`
	Amount       string `json:"amount"`
	From         string `json:"from,omitempty"`
	To           string `json:"to"`
}

type Tx struct {
	Block     TxBlockDetails `json:"block"`
	From      string         `json:"from"`
	To        string         `json:"to"`
	Amount    string         `json:"amount,omitempty"`
	Token     TxToken        `json:"token,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}
