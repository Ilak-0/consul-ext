package app

import (
	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func init() {
	app = gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
}

func Get() *gin.Engine {
	return app
}
