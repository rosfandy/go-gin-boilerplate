package config

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresDB *gorm.DB

func NewPostgresDB() (*gorm.DB, error) {
	return NewPostgresDBWithSSL(false)
}

func NewPostgresDBWithSSL(sslEnabled bool) (*gorm.DB, error) {
	if AppConfig == nil {
		return nil, fmt.Errorf("app config is not loaded")
	}

	pg := AppConfig.Postgres

	connURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(pg.User, pg.Password),
		Host:   fmt.Sprintf("%s:%s", pg.Host, pg.Port),
		Path:   pg.Database,
	}

	q := connURL.Query()
	if sslEnabled {
		q.Set("sslmode", "require")
	} else {
		q.Set("sslmode", "disable")
	}
	connURL.RawQuery = q.Encode()

	dsn := connURL.String()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitPostgresDB() (*gorm.DB, error) {
	return InitPostgresDBWithSSL(false)
}

func InitPostgresDBWithSSL(sslEnabled bool) (*gorm.DB, error) {
	if PostgresDB != nil {
		if AppConfig != nil && AppConfig.Postgres.Database != "" {
			currentDB, err := currentDatabase(PostgresDB)
			if err == nil && currentDB == AppConfig.Postgres.Database {
				return PostgresDB, nil
			}

			_ = ClosePostgresDB()
		} else {
			return PostgresDB, nil
		}
	}

	db, err := NewPostgresDBWithSSL(sslEnabled)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	if AppConfig != nil && AppConfig.Postgres.Database != "" {
		currentDB, err := currentDatabase(db)
		if err != nil {
			return nil, err
		}

		if currentDB != AppConfig.Postgres.Database {
			return nil, fmt.Errorf("connected to '%s' but expected '%s'", currentDB, AppConfig.Postgres.Database)
		}
	}

	PostgresDB = db
	return PostgresDB, nil
}

func ClosePostgresDB() error {
	if PostgresDB == nil {
		return nil
	}

	sqlDB, err := PostgresDB.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	PostgresDB = nil
	return nil
}

func currentDatabase(db *gorm.DB) (string, error) {
	var currentDB string
	if err := db.Raw("SELECT current_database()").Scan(&currentDB).Error; err != nil {
		return "", err
	}

	return currentDB, nil
}
