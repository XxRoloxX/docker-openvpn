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

	clientRouter := params.RootRouter.GinRouter.Group("/client")

	clientRouter.POST("", params.Handler.CreateClient)

	return &ClientRouter{}
}
