package client

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	PKI_PASSWORD_KEY = "PKI_PASSWORD"
)

type ClientService struct {
	pkiPassword     string
	clientDataStore ClientDataStore
	logger          *zap.Logger
}

type ClientServiceParams struct {
	fx.In
	ClientDataStore ClientDataStore
	Logger          *zap.Logger
}

func NewClientService(params *ClientServiceParams) *ClientService {

	pkiPasswd, isPkiPasswdSet := os.LookupEnv(PKI_PASSWORD_KEY)
	if !isPkiPasswdSet {
		panic(fmt.Sprintf("%s is not set", PKI_PASSWORD_KEY))
	}

	return &ClientService{
		pkiPassword:     pkiPasswd,
		logger:          params.Logger,
		clientDataStore: params.ClientDataStore,
	}
}

type NewClientData struct {
	ClientData *ClientData `json:"clientData"`
	Manifest   string      `json:"manifest"`
}

func (s *ClientService) CreateClient(clientName string, email string) (*NewClientData, error) {

	cmd := exec.Command("easyrsa",
		"--batch",
		fmt.Sprintf("--passin=pass:%s", s.pkiPassword),
		fmt.Sprintf("--passout=pass:%s", s.pkiPassword),
		"build-client-full",
		clientName,
		"nopass")

	exists, err := s.clientDataStore.DoesClientExists(clientName)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to check if client exists %s: %s", clientName, err)
		return nil, err
	}

	if exists {
		s.logger.Sugar().Errorf("Failed to client, because one with the same already exists %s: %s", clientName, err)
		return nil, errors.New(fmt.Sprintf("Client %s already exists", clientName))
	}

	_, err = cmd.Output()
	if err != nil {
		s.logger.Sugar().Errorf("Failed to create client with easy-rsa: %s", err)
		return nil, err
	}

	cmd = exec.Command("ovpn_getclient", clientName)
	manifest, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to get client .ovpn manifest", err)
		return nil, err
	}

	newClientData, err := s.clientDataStore.CreateClient(clientName, email)
	if err != nil {
		fmt.Println("Failed to create client in the db", err)
		return nil, err
	}

	return &NewClientData{
		ClientData: newClientData,
		Manifest:   string(manifest),
	}, nil
}

func (s *ClientService) RemoveClient(clientName string) error {

	cmd := exec.Command("ovpn_revokeclient", clientName)

	exists, err := s.clientDataStore.DoesClientExists(clientName)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to check if client exists %s: %s", clientName, err)
		return err
	}

	if !exists {
		s.logger.Sugar().Errorf("Failed to remove not-existing client %s: %s", clientName, err)
		return errors.New(fmt.Sprintf("Client %s does not exist", clientName))
	}

	_, err = cmd.Output()
	if err != nil {
		s.logger.Sugar().Errorf("Failed to revoke clients access %s", err)
		return err
	}

	err = s.clientDataStore.RemoveClient(clientName)
	if err != nil {
		s.logger.Sugar().Errorf("Failed to remove client %s from db: %s", clientName, err)
		return err
	}

	return nil
}

func (s *ClientService) GetClient(clientName string) (*ClientData, error) {

	clientData, err := s.clientDataStore.GetClient(clientName)
	if err != nil {
		s.logger.Sugar().Errorf("Client doesn't exist %s: %s", clientName, err)
		return nil, err
	}

	return clientData, nil
}

func (s *ClientService) GetClients() ([]*ClientData, error) {

	clientData, err := s.clientDataStore.GetAllClients()
	if err != nil {
		s.logger.Sugar().Errorf("Failed to get clients: %s", err)
		return nil, err
	}

	return clientData, nil
}

func (s *ClientService) DoesClientExists(clientName string) (bool, error) {
	return s.clientDataStore.DoesClientExists(clientName)
}
