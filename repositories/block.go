package repositories

import (
	"golang-block-chain/entities"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type BlockFactory func(ctx *GormDBContext) IBlock

type IBlock interface {
	ListBlocks(limit int) (blocks []*entities.Block, err error)
	Get(blockNum uint64) (block *entities.Block, err error)
	Save(block *entities.Block) (err error)
	//Delete(block *entities.Block) (err error)
}

type Block struct {
	GormRepository
}

func NewBlock(ctx *GormDBContext) IBlock {
	repository := new(Block)
	repository.SetDBContext(ctx)
	return repository
}

func (repo *Block) ListBlocks(limit int) ([]*entities.Block, error) {
	var blocks []*entities.Block
	err := repo.DB().Limit(limit).Order("block_time").Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (repo *Block) Get(blockNum uint64) (*entities.Block, error) {
	var block entities.Block
	err := repo.DB().Preload(clause.Associations).First(&block, "block_num = ?", blockNum).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &block, err
}

func (repo *Block) Save(block *entities.Block) (err error) {
	err = repo.DB().Save(block).Error
	if err != nil && isCausedByUniqueConstraint(err) {
		err = NewUniqueConstrainError(err)
	}
	return
}

//func (repo *Block) Delete(block *entities.Block) (err error) {
//	err = repo.DB().Delete(block).Error
//	return
//}
