package client

import (
	"managed-openvpn/internal/root-router"

	"go.uber.org/fx"
)

type ClientRouter struct {
	handler *ClientHandler
}

type ClientRouterParams struct {
	fx.In
	RootRouter *rootrouter.RootRouter
	Handler    *ClientHandler
}

func NewClientRouter(params ClientRouterParams) *ClientRouter {

	clientRouter := params.RootRouter.GinRouter.Group("/clients")

	clientRouter.POST("", params.Handler.CreateClient)
	clientRouter.GET("/:name", params.Handler.GetClient)
	clientRouter.GET("", params.Handler.GetClients)
	clientRouter.DELETE("/:name", params.Handler.RemoveClient)

	return &ClientRouter{}
}
