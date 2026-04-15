package db

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func ModelCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "model",
		Short: "Generate database model",
		RunE: func(cmd *cobra.Command, args []string) error {
			structName := toPascalCase(name)
			fileName := strings.ToLower(nonAlphaNum.ReplaceAllString(name, "_"))
			fileName = strings.Trim(fileName, "_")
			if fileName == "" {
				return fmt.Errorf("invalid model name: %s", name)
			}

			modelPath := filepath.Join("internal", "database", "model", fileName+".go")
			if _, err := os.Stat(modelPath); err == nil {
				return fmt.Errorf("model file already exists: %s", modelPath)
			} else if !os.IsNotExist(err) {
				return err
			}

			content := fmt.Sprintf("package model\n\n"+
				"type %s struct {\n"+
				"\tID uint `gorm:\"primaryKey\"`\n"+
				"}\n", structName)

			if err := os.WriteFile(modelPath, []byte(content), 0644); err != nil {
				return err
			}

			fmt.Println("Model", structName, "is generated at", modelPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Model name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func toPascalCase(input string) string {
	parts := nonAlphaNum.Split(input, -1)
	var b strings.Builder

	for _, part := range parts {
		if part == "" {
			continue
		}

		lower := strings.ToLower(part)
		b.WriteString(strings.ToUpper(lower[:1]))
		if len(lower) > 1 {
			b.WriteString(lower[1:])
		}
	}

	return b.String()
}
