package migration0

type Block struct {
	BlockNum     int64 `gorm:"primaryKey"`
	BlockHash    string
	BlockTime    uint64
	ParentHash   string
	Transactions []*Transaction `gorm:"foreignKey:BlockNum"`
}

type Transaction struct {
	TxHash   string `gorm:"primaryKey"`
	BlockNum int64
	From     string
	To       string
	Nonce    uint64
	Data     []byte
	Value    string
	Logs     []*TransactionLog `gorm:"foreignKey:TxHash"`
}

type TransactionLog struct {
	ID     int64 `gorm:"primaryKey"`
	TxHash string
	Index  uint
	Data   []byte
}
