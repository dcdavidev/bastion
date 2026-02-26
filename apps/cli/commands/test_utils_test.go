package commands

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = w
	os.Stderr = w

	// Disabilita colori e banner
	pterm.DisableColor()
	pterm.SetDefaultOutput(w)
	os.Setenv("BASTION_TEST", "true")
	defer os.Unsetenv("BASTION_TEST")

	root.SetArgs(args)
	root.SetOut(w)
	root.SetErr(w)

	err = root.Execute()
	
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output = buf.String()
	
	return output, err
}

func setupTestConfig(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "bastion-test-*")
	if err != nil {
		t.Fatal(err)
	}
	customConfigDir = tempDir
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
		customConfigDir = ""
	})
	return tempDir
}

func mockToken(t *testing.T, dir string, token string) {
	t.Helper()
	tokenPath := filepath.Join(dir, "token")
	err := os.WriteFile(tokenPath, []byte(token), 0600)
	if err != nil {
		t.Fatal(err)
	}
}
