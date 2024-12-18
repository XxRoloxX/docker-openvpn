package client

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	PKI_PASSWORD_KEY = "PKI_PASSWORD"
)

type ClientService struct {
	pkiPassword string
}

func NewClientService() *ClientService {

	pkiPasswd, isPkiPasswdSet := os.LookupEnv(PKI_PASSWORD_KEY)
	if !isPkiPasswdSet {
		panic(fmt.Sprintf("%s is not set", PKI_PASSWORD_KEY))
	}

	return &ClientService{
		pkiPassword: pkiPasswd,
	}
}

func (s *ClientService) CreateClient(clientName string) (string, error) {

	cmd := exec.Command("easy-rsa",
		"--batch",
		fmt.Sprintf("--passin=pass:%s", s.pkiPassword),
		fmt.Sprintf("--passout=pass:%s", s.pkiPassword),
		"build-client-full",
		clientName,
		"nopass")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to create client", err)
		return "", err
	}

	return string(output), nil
}
