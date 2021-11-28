//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"golang-block-chain/repositories"
	"golang-block-chain/services"
	"golang-block-chain/services/contexts"
	"golang-block-chain/web/controllers"
	"gorm.io/gorm"
)

func InitApp(
	db *gorm.DB,
	webEngine *gin.Engine,
	ethClient *ethclient.Client,
	startFrom int64,
) *App {
	wire.Build(
		wire.Struct(new(App), "*"),
		wire.Struct(new(Web), "*"),
		NewSyncer,
		wire.Value(repositories.BlockFactory(repositories.NewBlock)),
		wire.Value(repositories.TransactionFactory(repositories.NewTransaction)),
		wire.Struct(new(contexts.BlockFactory), "*"),
		wire.Struct(new(services.Block), "*"),
		wire.Bind(new(services.IBlock), new(*services.Block)),
		wire.Struct(new(controllers.Block), "*"),
		wire.Struct(new(controllers.Transaction), "*"),
	)
	return &App{}
}
