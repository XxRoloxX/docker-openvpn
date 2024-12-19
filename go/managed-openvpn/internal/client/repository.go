package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const BBOLT_CLIENT_STORE_PATH_KEY = "BBOLT_CLIENT_STORE_PATH"
const CLIENTS_BUCKET = "clients"

type ClientData struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type ClientDataStore interface {
	CreateClient(clientName string, email string) (*ClientData, error)
	RemoveClient(clientName string) error
	GetAllClients() ([]*ClientData, error)
	GetClient(clientName string) (*ClientData, error)
	DoesClientExists(clientName string) (bool, error)
}

type BboltClientDataStore struct {
	logger *zap.Logger
	db     *bbolt.DB
}

func NewBboltClientDataStore(logger *zap.Logger) *BboltClientDataStore {

	path, isPathSet := os.LookupEnv(BBOLT_CLIENT_STORE_PATH_KEY)
	if !isPathSet {
		panic(fmt.Sprintf("%s is not set", BBOLT_CLIENT_STORE_PATH_KEY))
	}

	bbolt, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		panic(
			fmt.Sprintf(
				"Failed to open the client store database at %s: %s",
				path,
				err.Error()),
		)
	}

	return &BboltClientDataStore{
		db:     bbolt,
		logger: logger,
	}
}

func (s *BboltClientDataStore) CreateClient(clientName string, email string) (*ClientData, error) {
	tx, err := s.db.Begin(true)

	if err != nil {
		s.logger.Error("Failed to start bbolt transaction", zap.Error(err))
		return nil, err
	}

	bucket, err := tx.CreateBucketIfNotExists([]byte(CLIENTS_BUCKET))

	currentClientData := bucket.Get([]byte(clientName))
	if currentClientData != nil {
		return nil,
			errors.New(fmt.Sprintf("Client: %s already exists in the database", clientName))
	}

	clientData := ClientData{
		Email:     email,
		Name:      clientName,
		UpdatedAt: time.Now().UnixMilli(),
		CreatedAt: time.Now().UnixMilli(),
	}

	encodedClient, err := json.Marshal(clientData)
	if err != nil {
		s.logger.Error("Failed to serialize client data", zap.Error(err))
		return nil, err
	}

	err = bucket.Put([]byte(clientName), encodedClient)

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit client data", zap.Error(err))
		return nil, err
	}

	return &clientData, nil
}

func (s *BboltClientDataStore) RemoveClient(clientName string) error {
	tx, err := s.db.Begin(true)

	if err != nil {
		s.logger.Error("Failed to start bbolt transaction", zap.Error(err))
		return err
	}

	exists, err := s.DoesClientExists(clientName)
	if err == nil {
		s.logger.Sugar().Errorf("Failed to check if client exists %s: %s", clientName, err)
		return err
	}

	if !exists {
		s.logger.Sugar().Errorf("Client doesn't exist: %s", clientName)
		return errors.New(fmt.Sprintf("Failed to remove non-existing client: %s", clientName))
	}

	bucket := tx.Bucket([]byte(CLIENTS_BUCKET))

	err = bucket.Delete([]byte(clientName))
	if err != nil {
		s.logger.Error("Failed to delete client data", zap.Error(err))
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit client data", zap.Error(err))
		return err
	}

	return nil
}

func (s *BboltClientDataStore) GetClient(clientName string) (*ClientData, error) {

	var clientData ClientData

	err := s.db.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket([]byte(CLIENTS_BUCKET))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Client: %s doesn't exists in the database", clientName))
		}

		currentClientData := bucket.Get([]byte(clientName))
		if currentClientData == nil {
			return errors.New(fmt.Sprintf("Client: %s doesn't exists in the database", clientName))
		}

		err := json.Unmarshal(currentClientData, &clientData)
		if err != nil {
			s.logger.Error("Failed to decode client data", zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to close a read-only transaction", zap.Error(err))
		return nil, err
	}

	return &clientData, nil
}

func (s *BboltClientDataStore) GetAllClients() ([]*ClientData, error) {

	clientData := make([]*ClientData, 0)

	err := s.db.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket([]byte(CLIENTS_BUCKET))
		if bucket == nil {
			return nil
		}

		bucket.ForEach(func(k, v []byte) error {
			var currnetClientData ClientData
			err := json.Unmarshal(v, &currnetClientData)
			if err != nil {
				s.logger.Error("Failed to decode client data", zap.Error(err))
				return err
			}

			clientData = append(clientData, &currnetClientData)
			return nil
		})

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to close a read-only transaction", zap.Error(err))
		return nil, err
	}

	return clientData, nil
}

func (s *BboltClientDataStore) DoesClientExists(clientName string) (bool, error) {

	var doesClientExists bool

	err := s.db.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket([]byte(CLIENTS_BUCKET))
		if bucket == nil {
			doesClientExists = false
			return nil
		}

		currentClientData := bucket.Get([]byte(clientName))
		if currentClientData == nil {
			doesClientExists = false
		} else {
			doesClientExists = true
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to close a read-only transaction", zap.Error(err))
		return false, err
	}

	return doesClientExists, nil
}

var _ ClientDataStore = &BboltClientDataStore{}
