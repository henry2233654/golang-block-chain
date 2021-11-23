package migration0

import (
	"time"
)

type Block struct {
	BlockNum     uint64 `gorm:"primaryKey"`
	BlockHash    string
	BlockTime    time.Time
	ParentHash   string
	Transactions []*Transaction `gorm:"foreignKey:BlockHash"`
}

type Transaction struct {
	TxHash    string `gorm:"primaryKey"`
	BlockHash string
	From      string
	To        string
	Nonce     uint64
	Data      string //input
	Value     string
	Logs      []*TransactionLog `gorm:"foreignKey:TxHash"`
}

type TransactionLog struct {
	TxHash string
	Index  uint64
	Data   string
}
