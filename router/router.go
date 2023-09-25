package router

import (
	_ "consul-ext/docs"
	_ "consul-ext/router/api/v1"
	"consul-ext/router/app"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	app.Get().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Run(addr ...string) error {
	return app.Get().Run(addr...)
}
