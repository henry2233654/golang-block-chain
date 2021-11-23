package entities

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
