package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/previder/vault-cli/pkg/model"
	"io"
	"log"
	"net/http"
)

const (
	DefaultBaseUri = "https://vault.previder.io"
	jsonEncoding   = "application/json; charset=utf-8"
)

type PreviderVaultClient interface {
	GetEnvironments() ([]model.Environment, error)
	GetEnvironment(id string) (*model.Environment, error)
	CreateEnvironment(create model.EnvironmentCreate) (*model.EnvironmentCreateResponse, error)
	DeleteEnvironment(id string) error

	GetTokens() ([]model.Token, error)
	GetToken(id string) (*model.Token, error)
	CreateToken(create model.TokenCreate) (*model.TokenCreateResponse, error)
	DeleteToken(id string) error

	GetSecrets() ([]model.Secret, error)
	GetSecret(id string) (*model.Secret, error)
	CreateSecret(create model.SecretCreate) (*model.Secret, error)
	DecryptSecret(id string) (*model.SecretDecrypt, error)
	DeleteSecret(id string) error
}

type VaultClient struct {
	PreviderVaultClient
	baseUri string
	token   string
	verbose bool
	http    *http.Client
}

func NewVaultClient(baseUri string, token string) (*VaultClient, error) {
	if baseUri == "" {
		baseUri = DefaultBaseUri
	}

	vaultClient := &VaultClient{
		baseUri: baseUri,
		token:   token,
	}

	if err := vaultClient.validateConnection(); err != nil {
		return nil, err
	}

	return vaultClient, nil
}

func (v *VaultClient) SetVerbose(verbose bool) {
	v.verbose = verbose
}

func (v *VaultClient) validateConnection() error {
	version := model.Version{Version: "test"}
	/*	err := v.request("GET", "/version", nil, version)
		if err != nil {
			return err
		}*/
	if v.verbose {
		log.Println(fmt.Sprintf("Setup new Vault client [%v] to %v", version.Version, v.baseUri))
	}
	return nil
}

func (v *VaultClient) request(method string, url string, requestBody interface{}, responseBody interface{}) error {

	// content will be empty with GET, so can be sent anyway
	if v == nil || v.baseUri == "" {
		return errors.New("Previder vault client not setup")
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, v.baseUri+url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", jsonEncoding)

	req.Header.Set("X-API-Key", v.token)

	req.Header.Set("Accept", jsonEncoding)

	httpClient := http.DefaultClient
	res, requestErr := httpClient.Do(req)
	if requestErr != nil {
		log.Printf("[ERROR] [Previder Vault] Error from API received")
		return requestErr
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if res.StatusCode == 401 {
			log.Fatal("Unauthorized")
		}
		log.Printf("An error was returned: %v, %v\n", res.StatusCode, res.Body)
	}

	if responseBody != nil {
		err := json.NewDecoder(res.Body).Decode(&responseBody)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

	return nil
}
