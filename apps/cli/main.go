package main

import (
	"fmt"
	"os"

	"github.com/dcdavidev/bastion/apps/cli/commands"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // Optional: ignore error if .env doesn't exist
	if err := commands.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
