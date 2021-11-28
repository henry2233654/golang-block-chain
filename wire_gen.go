// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"golang-block-chain/repositories"
	"golang-block-chain/services"
	"golang-block-chain/services/contexts"
	"golang-block-chain/web/controllers"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitApp(db *gorm.DB, webEngine *gin.Engine, ethClient *ethclient.Client, startFrom int64) *App {
	block := &services.Block{}
	blockFactory := _wireBlockFactoryValue
	transactionFactory := _wireTransactionFactoryValue
	contextsBlockFactory := &contexts.BlockFactory{
		DB:                     db,
		BlockRepoFactory:       blockFactory,
		TransactionRepoFactory: transactionFactory,
	}
	controllersBlock := &controllers.Block{
		Service:    block,
		CtxFactory: contextsBlockFactory,
	}
	transaction := &controllers.Transaction{
		Service:    block,
		CtxFactory: contextsBlockFactory,
	}
	web := Web{
		Block:       controllersBlock,
		Transaction: transaction,
	}
	blockFactory2 := contexts.BlockFactory{
		DB:                     db,
		BlockRepoFactory:       blockFactory,
		TransactionRepoFactory: transactionFactory,
	}
	syncer := NewSyncer(ethClient, block, blockFactory2, startFrom)
	app := &App{
		Web:         web,
		DB:          db,
		WebEngine:   webEngine,
		ChainSyncer: syncer,
	}
	return app
}

var (
	_wireBlockFactoryValue       = repositories.BlockFactory(repositories.NewBlock)
	_wireTransactionFactoryValue = repositories.TransactionFactory(repositories.NewTransaction)
)
