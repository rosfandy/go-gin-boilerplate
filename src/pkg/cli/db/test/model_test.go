package test

import (
	"strings"
	"testing"

	"owner-api-proxy/pkg/cli/db"
)

func TestModelCommand_HasExpectedDefaultFlag(t *testing.T) {
	cmd := db.ModelCommand()

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		t.Fatalf("failed to read name flag: %v", err)
	}
	if name != "" {
		t.Fatalf("expected default name to be empty, got %q", name)
	}
}

func TestModelCommand_RequiresNameFlag(t *testing.T) {
	cmd := db.ModelCommand()

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when --name is not provided")
	}

	if !strings.Contains(err.Error(), "required flag(s) \"name\" not set") {
		t.Fatalf("unexpected error message: %q", err.Error())
	}
}

func TestModelCommand_ReturnsErrorForInvalidName(t *testing.T) {
	cmd := db.ModelCommand()
	cmd.SetArgs([]string{"--name", "!!!"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid model name")
	}

	want := "invalid model name: !!!"
	if err.Error() != want {
		t.Fatalf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}
