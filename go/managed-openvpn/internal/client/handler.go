package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ClientHandler struct {
	service *ClientService
	logger  *zap.Logger
}

type ClientHandlerParams struct {
	fx.In
	Service *ClientService
	Logger  *zap.Logger
}

func NewClientHandler(params ClientHandlerParams) *ClientHandler {
	return &ClientHandler{
		service: params.Service,
		logger:  params.Logger,
	}
}

type CreateClientParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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

	newClient, err := h.service.CreateClient(params.Name, params.Email)
	if err != nil {
		g.JSON(http.StatusInternalServerError, "Failed to create client")
		return
	}

	g.JSON(
		http.StatusCreated,
		newClient,
	)
}

func (h *ClientHandler) RemoveClient(g *gin.Context) {

	name := g.Param("name")

	if name == "" {
		g.JSON(
			http.StatusBadRequest,
			"Missing 'name' parameter",
		)
		return
	}

	exists, err := h.service.DoesClientExists(name)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err)
		return
	}

	if !exists {
		g.JSON(http.StatusBadRequest, fmt.Sprintf("Client %s doesn't exist", name))
		return
	}

	err = h.service.RemoveClient(name)
	if err != nil {
		g.JSON(http.StatusInternalServerError, "Failed to remove client")
		return
	}

	g.JSON(
		http.StatusOK,
		nil,
	)
}

func (h *ClientHandler) GetClient(g *gin.Context) {

	name := g.Param("name")

	if name == "" {
		g.JSON(
			http.StatusBadRequest,
			"Missing 'name' parameter",
		)
		return
	}

	exists, err := h.service.DoesClientExists(name)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err)
		return
	}

	if !exists {
		g.JSON(http.StatusNotFound, fmt.Sprintf("Client %s doesn't exist", name))
		return
	}

	clientData, err := h.service.GetClient(name)
	if err != nil {
		g.JSON(http.StatusInternalServerError, "Failed to get client")
		return
	}

	g.JSON(
		http.StatusOK,
		clientData,
	)
}

func (h *ClientHandler) GetClients(g *gin.Context) {

	clients, err := h.service.GetClients()
	if err != nil {
		g.JSON(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, clients)
}
