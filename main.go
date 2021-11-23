package main

import (
	"golang-block-chain/configs"

	"github.com/gin-gonic/gin"
)

var AppGinEngine *gin.Engine
var Config configs.Config

// @Title System Main Data
func main() {
	config := configs.GetConfig()
	_, db := InitializeDb(config.RelationalDB)
	defer CloseDb(db)
	webEngine := gin.Default()
	app := InitApp(db, webEngine)
	_ = app.Serve()
}
