package config

import (
	"fmt"
	"net/url"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SqlInstance struct {
	connection string // "postgres" or "mysql"
	client     *gorm.DB
}

func NewClientSql(connection string) *SqlInstance {
	instance := SqlInstance{
		connection: connection,
	}
	return &instance
}

func (sql *SqlInstance) Open() error {
	if sql.connection == "" {
		return fmt.Errorf("connection type is required")
	}

	if sql.client != nil {
		return nil
	}

	if AppConfig == nil {
		return fmt.Errorf("app config is not loaded")
	}

	connection := strings.ToLower(sql.connection)
	var (
		db  *gorm.DB
		err error
	)

	switch connection {
	case "postgres":
		dsn := buildPostgresDSN(AppConfig.Postgres)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "mysql":
		dsn := buildMySQLDSN(AppConfig.MySQL)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		return fmt.Errorf("unsupported connection type: %s", sql.connection)
	}

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	sql.client = db
	return nil
}

func (sql *SqlInstance) Close() error {
	if sql.client == nil {
		return nil
	}

	sqlDB, err := sql.client.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	sql.client = nil

	return nil
}

func (sql *SqlInstance) Client() *gorm.DB {
	return sql.client
}

func buildPostgresDSN(cfg DbConfig) string {
	sslMode := cfg.Secure
	if sslMode == "" {
		sslMode = "disable"
	}

	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.Database,
	}

	q := connURL.Query()
	q.Set("sslmode", sslMode)
	connURL.RawQuery = q.Encode()

	return connURL.String()
}

func buildMySQLDSN(cfg DbConfig) string {
	params := "charset=utf8mb4&parseTime=True&loc=Local"
	if cfg.Secure != "" {
		tlsMode := strings.ToLower(strings.TrimSpace(cfg.Secure))
		switch tlsMode {
		case "enable":
			tlsMode = "true"
		case "disable":
			tlsMode = "false"
		}

		params = params + "&tls=" + url.QueryEscape(tlsMode)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, params)
}
