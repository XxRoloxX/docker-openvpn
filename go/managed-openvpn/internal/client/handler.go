package client

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type ClientHandler struct {
	service *ClientService
}

type ClientHandlerParams struct {
	fx.In
	Service *ClientService
}

func NewClientHandler(params *ClientHandlerParams) *ClientHandler {
	return &ClientHandler{
		service: params.Service,
	}
}

func (h *ClientHandler) CreateClient(g *gin.Context) {

	exec.Command("easy-rsa", "build-client-full", "client-name", "nopass")

	g.JSON(http.StatusOK, "The client endpoint works")
}
