package services

import (
	"golang-block-chain/entities"
	"golang-block-chain/services/contexts"
)

type IBlock interface {
	ListBlocks(c *contexts.Block, limit int) ([]*entities.Block, error)
	RetrieveBlock(c *contexts.Block, blockNum int64) (*entities.Block, error)
	RetrieveTransaction(c *contexts.Block, txHash string) (*entities.Transaction, error)
	SaveBlock(c *contexts.Block, block *entities.Block) error
}

type Block struct {
}

func (srv *Block) ListBlocks(c *contexts.Block, limit int) (blocks []*entities.Block, err error) {
	return c.BlockRepo.ListBlocks(limit)
}

func (srv *Block) RetrieveBlock(c *contexts.Block, blockNum int64) (*entities.Block, error) {
	block, err := c.BlockRepo.Get(blockNum)
	if err != nil {
		return nil, err
	}
	if block == nil {
		err := &NotExistError{ResourceName: "block", Expected: blockNum}
		return nil, err
	}
	return block, nil
}

func (srv *Block) RetrieveTransaction(c *contexts.Block, txHash string) (*entities.Transaction, error) {
	transaction, err := c.TransactionRepo.Get(txHash)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		err := &NotExistError{ResourceName: "transaction", Expected: txHash}
		return nil, err
	}
	return transaction, nil
}

func (srv *Block) SaveBlock(c *contexts.Block, block *entities.Block) error {
	return StartTransaction(c, func() error {
		err := c.BlockRepo.Save(block)
		return err
	})
}
