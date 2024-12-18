package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"managed-openvpn/internal/client"
	rootrouter "managed-openvpn/internal/root-router"
	"net/http"
)

type ApplicationParams struct {
	fx.In
	RootRouter   *rootrouter.RootRouter
	ClientRouter *client.ClientRouter
}

type Application struct {
	rootRouter   *rootrouter.RootRouter
	clientRouter *client.ClientRouter
}

func NewApplication(lc fx.Lifecycle, params ApplicationParams) *Application {

	srv := &http.Server{
		Addr:    ":8080",
		Handler: params.RootRouter.GinRouter.Handler(),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return &Application{
		rootRouter:   params.RootRouter,
		clientRouter: params.ClientRouter,
	}
}

func main() {
	fx.New(
		fx.Provide(
			rootrouter.NewRootRouter,
			client.NewClientRouter,
			client.NewClientHandler,
			client.NewClientService,
			NewApplication,
		),
		fx.Invoke(func(*Application) {}),
	).Run()
}
