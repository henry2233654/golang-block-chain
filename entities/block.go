package entities

import "time"

type Block struct {
	BlockNum     uint64 `gorm:"primaryKey"`
	BlockHash    string
	BlockTime    time.Time
	ParentHash   string
	Transactions []*Transaction `gorm:"foreignKey:BlockHash"`
}
