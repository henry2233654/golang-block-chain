package entities

type Block struct {
	BlockNum     int64 `gorm:"primaryKey"`
	BlockHash    string
	BlockTime    uint64
	ParentHash   string
	Transactions []*Transaction `gorm:"foreignKey:BlockNum"`
}
