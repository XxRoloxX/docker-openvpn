package rootrouter

import (
	"github.com/gin-gonic/gin"
)

type RootRouter struct {
	GinRouter *gin.Engine
}

func NewRootRouter() *RootRouter {
	return &RootRouter{
		GinRouter: gin.Default(),
	}
}
