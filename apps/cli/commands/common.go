package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dcdavidev/bastion/packages/core/models"
)

func loadToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(home, ".bastion", "token")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type vaultConfigResponse struct {
	WrappedMasterKey string `json:"wrapped_master_key"`
	MasterKeySalt    string `json:"master_key_salt"`
}

func fetchVaultConfig(url, token string) (*vaultConfigResponse, error) {
	req, _ := http.NewRequest("GET", url+"/api/v1/vault/config", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch vault config: %s", resp.Status)
	}

	var config vaultConfigResponse
	json.NewDecoder(resp.Body).Decode(&config)
	return &config, nil
}

func fetchProject(url, token, id string) (*models.Project, error) {
	req, _ := http.NewRequest("GET", url+"/api/v1/projects/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch project: %s", resp.Status)
	}

	var project models.Project
	json.NewDecoder(resp.Body).Decode(&project)
	return &project, nil
}
