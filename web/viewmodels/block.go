package viewmodels

import (
	"time"
)

type ListBlocks struct {
	Blocks []ListBlockItem `json:"blocks"`
}

type ListBlockItem struct {
	BlockNum   uint64    `json:"block_num"`
	BlockHash  string    `json:"block_hash"`
	BlockTime  time.Time `json:"block_time"`
	ParentHash string    `json:"parent_hash"`
}

type GetSingleBlock struct {
	BlockNum     uint64    `json:"block_num"`
	BlockHash    string    `json:"block_hash"`
	BlockTime    time.Time `json:"block_time"`
	ParentHash   string    `json:"parent_hash"`
	Transactions []string  `json:"transactions"`
}
