package rootrouter

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type RootRouter struct {
	GinRouter *gin.Engine
}

func NewRootRouter() *RootRouter {

	root := gin.Default()
	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &RootRouter{
		GinRouter: root,
	}
}
