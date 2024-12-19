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

// ClientHandler handles client-related API calls.
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

// CreateClient godoc
// @Summary Create a new client
// @Description Create a new client with the given name and email
// @Tags clients
// @Accept json
// @Produce json
// @Param client body CreateClientParams true "Client details"
// @Success 201 {object} NewClientData
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /clients [post]
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

	if params.Name == "" {
		g.JSON(http.StatusBadRequest, "name parameter cannot be empty")
		return
	}

	if params.Email == "" {
		g.JSON(http.StatusBadRequest, "email parameter cannot be empty")
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

// RemoveClient godoc
// @Summary Remove a client
// @Description Remove a client by name
// @Tags clients
// @Produce json
// @Param name path string true "Client name"
// @Success 204 {null} null "Client removed successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /clients/{name} [delete]
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

	g.Status(
		http.StatusNoContent,
	)
}

// GetClient godoc
// @Summary Get a client by name
// @Description Get details of a client by their name
// @Tags clients
// @Produce json
// @Param name path string true "Client name"
// @Success 200 {object} ClientData
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Client not found"
// @Failure 500 {string} string "Internal server error"
// @Router /clients/{name} [get]
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

// GetClients godoc
// @Summary Get all clients
// @Description Retrieve a list of all clients
// @Tags clients
// @Produce json
// @Success 200 {array} ClientData
// @Failure 500 {string} string "Internal server error"
// @Router /clients [get]
func (h *ClientHandler) GetClients(g *gin.Context) {
	clients, err := h.service.GetClients()
	if err != nil {
		g.JSON(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, clients)
}
