package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"owner-api-proxy/internal/config"
	"owner-api-proxy/internal/database/migration"
	"owner-api-proxy/pkg/cli/db"

	"github.com/spf13/cobra"
)

func DbCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "db <command>",
		Short:   "Database command",
		Long:    "Command for database utilites",
		Example: "db pull",
	}

	cmd.AddCommand(
		db.PullCommand(runPull),
		db.MigrateCommand(runMigrate),
		db.ModelCommand(),
		db.PingCommand(runPing),
	)

	return cmd
}

func runMigrate(options db.MigrateOptions) error {
	if options.Down {
		return migration.DownWithDB(&options.ConfigPath, options.SSLMode, options.DBType)
	}

	return migration.UpWithDB(&options.ConfigPath, options.SSLMode, options.DBType)
}

func runPing(options db.PingOptions) error {
	if _, err := config.LoadConfig(&options.ConfigPath); err != nil {
		return err
	}

	sqlClient := config.NewClientSql(options.DBType)
	if err := sqlClient.Open(); err != nil {
		return err
	}
	defer sqlClient.Close()

	dbClient := sqlClient.Client()
	if dbClient == nil {
		return fmt.Errorf("database client is not initialized")
	}

	sqlDB, err := dbClient.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	fmt.Println("database connection is successful for", options.DBType)
	return nil
}

func runPull(options db.PullOptions) error {
	if _, err := config.LoadConfig(&options.ConfigPath); err != nil {
		return err
	}

	sqlClient := config.NewClientSql(options.DBType)
	if err := sqlClient.Open(); err != nil {
		return err
	}
	defer sqlClient.Close()

	dbClient := sqlClient.Client()
	if dbClient == nil {
		return fmt.Errorf("database client is not initialized")
	}

	columns, err := dbClient.Migrator().ColumnTypes(options.Name)
	if err != nil {
		return err
	}

	if len(columns) == 0 {
		return fmt.Errorf("table not found or has no columns: %s", options.Name)
	}

	fileName := strings.ToLower(db.NonAlphaNum().ReplaceAllString(options.Name, "_"))
	fileName = strings.Trim(fileName, "_")
	if fileName == "" {
		return fmt.Errorf("invalid table name: %s", options.Name)
	}

	modelPath := filepath.Join("internal", "database", "model", fileName+".go")
	if _, err := os.Stat(modelPath); err == nil {
		return fmt.Errorf("model file already exists: %s", modelPath)
	} else if !os.IsNotExist(err) {
		return err
	}

	structName := db.ToPascalCase(options.Name)
	var builder strings.Builder
	builder.WriteString("package model\n\n")
	builder.WriteString("type ")
	builder.WriteString(structName)
	builder.WriteString(" struct {\n")

	for _, col := range columns {
		columnName := col.Name()
		nullable, ok := col.Nullable()
		if !ok {
			nullable = true
		}

		fieldName := db.ToPascalCase(columnName)
		fieldType := db.PostgresTypeToGoType(col.DatabaseTypeName(), nullable)
		isPrimaryKey, _ := col.PrimaryKey()

		builder.WriteString("\t")
		builder.WriteString(fieldName)
		builder.WriteString(" ")
		builder.WriteString(fieldType)
		builder.WriteString(" `gorm:\"column:")
		builder.WriteString(columnName)
		if isPrimaryKey || strings.EqualFold(columnName, "id") {
			builder.WriteString(";primaryKey")
		}
		builder.WriteString("\" json:\"")
		builder.WriteString(columnName)
		builder.WriteString("\"`\n")
	}

	builder.WriteString("}\n\n")
	builder.WriteString("func (")
	builder.WriteString(structName)
	builder.WriteString(") TableName() string {\n")
	builder.WriteString("\treturn \"")
	builder.WriteString(options.Name)
	builder.WriteString("\"\n}\n")

	if err := os.WriteFile(modelPath, []byte(builder.String()), 0644); err != nil {
		return err
	}

	fmt.Println("Model", structName, "is generated at", modelPath)
	return nil
}
