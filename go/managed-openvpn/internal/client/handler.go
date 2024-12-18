package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func NewClientHandler(params ClientHandlerParams) *ClientHandler {
	return &ClientHandler{
		service: params.Service,
	}
}

type CreateClientParams struct {
	Name string `json:"name"`
}

func (h *ClientHandler) CreateClient(g *gin.Context) {

	body, err := io.ReadAll(g.Request.Body)
	if err != nil {
		g.JSON(http.StatusInternalServerError, "Failed to read request body")
		return
	}

	params := CreateClientParams{}

	err = json.Unmarshal(body, &params)
	if err != nil {
		g.JSON(http.StatusInternalServerError, "Failed to decode body")
		return
	}

	output, err := h.service.CreateClient(params.Name)

	g.JSON(
		http.StatusOK,
		fmt.Sprintf("The client endpoint works: %s", output),
	)
}
