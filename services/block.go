package services

import (
	"golang-block-chain/entities"
	"golang-block-chain/services/contexts"
)

//type CreateBlock struct {
//	No        string `json:"no"`         // 文章代號(unique)
//	Name      string `json:"name"`       // 文章名稱
//	QuickCode string `json:"quick_code"` // 簡碼
//}
//
//type UpdateBlock struct {
//	ID        uint   `json:"id"`         // 流水號
//	No        string `json:"no"`         // 文章代號(unique)
//	Name      string `json:"name"`       // 文章名稱
//	QuickCode string `json:"quick_code"` // 簡碼
//}

type IBlock interface {
	ListBlocks(c *contexts.Block, limit int) ([]*entities.Block, error)

	//CreateBlock(c *contexts.Block, index int, createBlock CreateBlock) (block *entities.Block, err error)
	//BatchCreateBlocks(c *contexts.Block, createBlocks []CreateBlock) (err error)
	//BatchUpdateBlocks(c *contexts.Block, updateBlocks []UpdateBlock) (err error)
	//BatchDeleteBlocks(c *contexts.Block, ids []uint) (err error)

	RetrieveBlock(c *contexts.Block, blockNum uint64) (*entities.Block, error)
	RetrieveTransaction(c *contexts.Block, txHash string) (*entities.Transaction, error)
}

type Block struct {
}

func (srv *Block) ListBlocks(c *contexts.Block, limit int) (blocks []*entities.Block, err error) {
	return c.BlockRepo.ListBlocks(limit)
}

func (srv *Block) RetrieveBlock(c *contexts.Block, blockNum uint64) (*entities.Block, error) {
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

//func (srv *Block) CreateBlock(c *contexts.Block, index int, createBlock CreateBlock) (block *entities.Block, err error) {
//	block = new(entities.Block)
//	block.No = createBlock.No
//	block.Name = createBlock.Name
//	block.QuickCode = createBlock.QuickCode
//	if err := srv.saveBlock(c, index, block); err != nil {
//		return nil, err
//	}
//	return block, nil
//}
//
//func (srv *Block) BatchCreateBlocks(c *contexts.Block, createBlocks []CreateBlock) (err error) {
//	fn := func() (err error) {
//		for i, createBlock := range createBlocks {
//			if _, err := srv.CreateBlock(c, i, createBlock); err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//	return StartTransaction(c, fn)
//}
//
//func (srv *Block) BatchUpdateBlocks(c *contexts.Block, updateBlocks []UpdateBlock) (err error) {
//	fn := func() (err error) {
//		for i, updateBlock := range updateBlocks {
//			id := updateBlock.ID
//			block, err := srv.RetrieveBlock(c, id)
//			if err != nil {
//				return err
//			}
//			block.No = updateBlock.No
//			block.Name = updateBlock.Name
//			block.QuickCode = updateBlock.QuickCode
//			if err := srv.saveBlock(c, i, block); err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//	return StartTransaction(c, fn)
//}

//func (srv *Block) BatchDeleteBlocks(c *contexts.Block, ids []uint) (err error) {
//	fn := func() (err error) {
//		for i, id := range ids {
//			block, err := srv.RetrieveBlock(c, id)
//			if err != nil {
//				return err
//			}
//			if err = srv.validateBeforeDelete(c, i, block); err != nil {
//				return err
//			}
//			if err = c.BlockRepo.Delete(block); err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//	return StartTransaction(c, fn)
//}

//func (srv *Block) saveBlock(c *contexts.Block, index int, block *entities.Block) error {
//	if err := srv.validateBeforeSave(c, block); err != nil {
//		return err
//	}
//	if err := c.BlockRepo.Save(block); err != nil {
//		if uniqueConstrainError, ok := err.(*repositories.UniqueConstrainError); ok {
//			return NewDuplicateError(&index, uniqueConstrainError)
//		}
//		return err
//	}
//	return nil
//}
