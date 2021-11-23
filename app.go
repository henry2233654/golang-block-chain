package main

import (
	"github.com/gin-gonic/gin"
	"golang-block-chain/migrations"
	_webControllers "golang-block-chain/web/controllers"
	"gorm.io/gorm"
	"reflect"
)

type App struct {
	Web
	DB        *gorm.DB
	WebEngine *gin.Engine
}

func (a *App) Serve() error {
	a.migrate()
	a.Web.register(a.WebEngine)
	return a.WebEngine.Run(":8081")
}

func (a *App) migrate() {
	migrations.Migrate(a.DB)
}

type Web struct {
	Block       *_webControllers.Block
	Transaction *_webControllers.Transaction
}

func (w Web) register(e *gin.Engine) {
	type controller interface {
		Route(*gin.Engine)
	}
	v := reflect.ValueOf(w)
	for i := 0; i < v.NumField(); i++ {
		c, ok := v.Field(i).Interface().(controller)
		if !ok {
			continue
		}
		c.Route(e)
	}
}
