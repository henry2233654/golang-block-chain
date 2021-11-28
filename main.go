package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"golang-block-chain/configs"
	"strconv"

	"github.com/gin-gonic/gin"
)

//@Title System Main Data
func main() {
	config := configs.GetConfig()
	_, db := InitializeDb(config.RelationalDB)
	defer CloseDb(db)
	webEngine := gin.Default()
	ethClient, err := ethclient.Dial(config.BlockChain.RpcUrl)
	if err != nil {
		panic(err)
	}
	from, err := strconv.ParseInt(config.BlockChain.From, 10, 64)
	if err != nil {
		panic(err)
	}
	app := InitApp(db, webEngine, ethClient, from)
	_ = app.Serve()
}
