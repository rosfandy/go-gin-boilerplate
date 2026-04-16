package test

import (
	"testing"

	"owner-api-proxy/pkg/cli/db"
)

func TestMigrateCommand_HasExpectedDefaultFlags(t *testing.T) {
	cmd := db.MigrateCommand()

	down, err := cmd.Flags().GetBool("down")
	if err != nil {
		t.Fatalf("failed to read down flag: %v", err)
	}
	if down {
		t.Fatalf("expected default down to be false")
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		t.Fatalf("failed to read config flag: %v", err)
	}
	if configPath != "app.yaml" {
		t.Fatalf("expected default config to be app.yaml, got %q", configPath)
	}

	ssl, err := cmd.Flags().GetString("ssl")
	if err != nil {
		t.Fatalf("failed to read ssl flag: %v", err)
	}
	if ssl != "" {
		t.Fatalf("expected default ssl to be empty, got %q", ssl)
	}
}

func TestMigrateCommand_ReturnsErrorForInvalidSSLValue(t *testing.T) {
	cmd := db.MigrateCommand()
	cmd.SetArgs([]string{"--ssl", "invalid"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid ssl value")
	}

	want := "invalid --ssl value: invalid (use enable or disable)"
	if err.Error() != want {
		t.Fatalf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}
