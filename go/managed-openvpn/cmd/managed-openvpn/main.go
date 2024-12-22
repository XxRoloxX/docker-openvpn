package main

import (
	"context"
	"fmt"
	"managed-openvpn/internal/client"
	rootrouter "managed-openvpn/internal/root-router"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
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
			zap.NewProduction,
			rootrouter.NewRootRouter,
			client.NewClientRouter,
			client.NewClientHandler,
			client.NewClientService,
			NewApplication,
			fx.Annotate(
				client.NewBboltClientDataStore,
				fx.As(new(client.ClientDataStore)),
			),
		),
		fx.Invoke(func(*Application) {}),
	).Run()
}
