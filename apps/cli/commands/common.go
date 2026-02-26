package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dcdavidev/bastion/packages/core/models"
	"github.com/dcdavidev/bastion/packages/core/version"
	"github.com/pterm/pterm"
)

var customConfigDir string

func getConfigDir() (string, error) {
        if customConfigDir != "" {
                return customConfigDir, nil
        }
        home, err := os.UserHomeDir()
        if err != nil {
                return "", err
        }
        return filepath.Join(home, ".bastion"), nil
}

func loadToken() (string, error) {
        configDir, err := getConfigDir()
        if err != nil {
                return "", err
        }
        configPath := filepath.Join(configDir, "token")
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

// CheckForUpdates checks GitHub for the latest release and displays a warning if a new version is available.
// It caches the last check time to avoid frequent API calls.
func CheckForUpdates() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	cachePath := filepath.Join(home, ".bastion", "last_update_check")
	
	// 1. Check if we checked recently (last 24 hours)
	if info, err := os.Stat(cachePath); err == nil {
		if time.Since(info.ModTime()) < 24*time.Hour {
			return
		}
	}

	// Create config dir if missing
	_ = os.MkdirAll(filepath.Dir(cachePath), 0700)

	// Update cache timestamp regardless of success to avoid spamming on network failure
	_ = os.WriteFile(cachePath, []byte(time.Now().String()), 0600)

	// 2. Fetch latest version from GitHub
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/dcdavidev/bastion/releases/latest")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return
	}

	// 3. Compare versions
	latest := strings.TrimPrefix(release.TagName, "v")
	current := strings.TrimPrefix(version.Version, "v")

	if isNewer(latest, current) {
		pterm.DefaultBox.
			WithTitle("Update Available").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgYellow)).
			Printf("A new version of Bastion CLI is available: %s (current: %s)\n\nPlease install it using: npm install -g @dcdavidev/bastion-cli", 
				pterm.Bold.Sprint(latest), pterm.FgGray.Sprint(current))
		fmt.Println()
	}
}

// isNewer performs a simple semver comparison.
func isNewer(latest, current string) bool {
	if latest == current {
		return false
	}
	
	lParts := strings.Split(latest, ".")
	cParts := strings.Split(current, ".")
	
	for i := 0; i < len(lParts) && i < len(cParts); i++ {
		if lParts[i] > cParts[i] {
			return true
		}
		if lParts[i] < cParts[i] {
			return false
		}
	}
	
	return len(lParts) > len(cParts)
}
