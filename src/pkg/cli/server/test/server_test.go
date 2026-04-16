package test

import (
	"errors"
	"strings"
	"testing"

	"owner-api-proxy/pkg/cli/server"
)

func TestServerCommand_HasExpectedDefaultFlags(t *testing.T) {
	cmd := server.ServerCommand(func(configPath string) error {
		return nil
	})

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		t.Fatalf("failed to read config flag: %v", err)
	}
	if configPath != "app.yaml" {
		t.Fatalf("expected default config to be app.yaml, got %q", configPath)
	}
}

func TestServerCommand_ReturnsErrorForMissingConfigFile(t *testing.T) {
	cmd := server.ServerCommand(func(configPath string) error {
		return errors.New("open " + configPath + ": no such file or directory")
	})
	cmd.SetArgs([]string{"--config", "missing-config.yaml"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing config file")
	}

	if !strings.Contains(err.Error(), "missing-config.yaml") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}
