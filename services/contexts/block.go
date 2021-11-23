package contexts

import (
	"golang-block-chain/repositories"
	"gorm.io/gorm"
)

type Block struct {
	Context
	BlockRepo       repositories.IBlock
	TransactionRepo repositories.ITransaction
	UserID          *uint
}

type BlockFactory struct {
	DB                     *gorm.DB
	BlockRepoFactory       repositories.BlockFactory
	TransactionRepoFactory repositories.TransactionFactory
}

func (f *BlockFactory) NewContext() *Block {
	dbCtx := repositories.NewGormDBContext(f.DB)
	ctx := new(Block)
	ctx.BlockRepo = f.BlockRepoFactory(dbCtx)
	ctx.TransactionRepo = f.TransactionRepoFactory(dbCtx)
	ctx.AddDBContexts(dbCtx)
	return ctx
}
