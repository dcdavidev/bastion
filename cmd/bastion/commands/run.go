package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var projectID string

var runCmd = &cobra.Command{
	Use:   "run -- [command]",
	Short: "Inject secrets into a command's environment",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverURL, _ := cmd.Flags().GetString("url")
		
		token, err := loadToken()
		if err != nil {
			return fmt.Errorf("not authenticated. Please run 'bastion login' first: %w", err)
		}

		// 1. Fetch Secrets (for now we just print them, E2EE logic will follow)
		secrets, err := fetchSecrets(serverURL, token, projectID)
		if err != nil {
			return err
		}

		// 2. Prepare environment
		env := os.Environ()
		for key, value := range secrets {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}

		// 3. Execute command
		externalCmd := exec.Command(args[0], args[1:]...)
		externalCmd.Env = env
		externalCmd.Stdout = os.Stdout
		externalCmd.Stderr = os.Stderr
		externalCmd.Stdin = os.Stdin

		return externalCmd.Run()
	},
}

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

func fetchSecrets(url, token, projectID string) (map[string]string, error) {
	req, _ := http.NewRequest("GET", url+"/api/v1/secrets?project_id="+projectID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch secrets: %s", resp.Status)
	}

	// This is a simplified placeholder. Real E2EE decryption will be added.
	var secretsData []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	json.NewDecoder(resp.Body).Decode(&secretsData)

	results := make(map[string]string)
	for _, s := range secretsData {
		results[s.Key] = s.Value
	}
	return results, nil
}

func init() {
	runCmd.Flags().StringVarP(&projectID, "project", "p", "", "Project ID to fetch secrets from")
	runCmd.MarkFlagRequired("project")
	rootCmd.AddCommand(runCmd)
}
