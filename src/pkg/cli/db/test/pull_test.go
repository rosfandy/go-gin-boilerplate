package test

import (
	"errors"
	"strings"
	"testing"

	"owner-api-proxy/pkg/cli/db"
)

func TestPullCommand_HasExpectedDefaultFlags(t *testing.T) {
	cmd := db.PullCommand(func(options db.PullOptions) error {
		return nil
	})

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		t.Fatalf("failed to read name flag: %v", err)
	}
	if name != "" {
		t.Fatalf("expected default name to be empty, got %q", name)
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		t.Fatalf("failed to read config flag: %v", err)
	}
	if configPath != "app.yaml" {
		t.Fatalf("expected default config to be app.yaml, got %q", configPath)
	}
}

func TestPullCommand_RequiresNameFlag(t *testing.T) {
	cmd := db.PullCommand(func(options db.PullOptions) error {
		return errors.New("should not be called")
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --name is not provided")
	}

	if !strings.Contains(err.Error(), "required flag(s) \"name\" not set") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}
