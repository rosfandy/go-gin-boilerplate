package migration

import (
	"fmt"
	"owner-api-proxy/internal/config"
	"owner-api-proxy/internal/database/model"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var migrationModels = []any{
	&model.Users{},
}

func Up(configPath *string, sslMode string) error {
	log := config.NewLogrusWithCategory("migration")
	log.Info("running up migration")

	if _, err := config.LoadConfig(configPath); err != nil {
		log.WithError(err).Error("failed to load config")
		return err
	}

	sslEnabled := sslMode == "enable"

	if config.AppConfig == nil || config.AppConfig.Postgres.Database == "" {
		err := fmt.Errorf("postgres database is not configured")
		log.WithError(err).Error("failed to resolve target database")
		return err
	}

	log.WithFields(map[string]any{
		"host":     config.AppConfig.Postgres.Host,
		"port":     config.AppConfig.Postgres.Port,
		"user":     config.AppConfig.Postgres.User,
		"database": config.AppConfig.Postgres.Database,
		"ssl":      sslEnabled,
	}).Info("target postgres connection")

	db, err := config.InitPostgresDBWithSSL(sslEnabled)
	if err != nil {
		log.WithError(err).Error("failed to initialize postgres")
		return err
	}
	defer config.ClosePostgresDB()

	var currentDatabase string
	if err := db.Raw("SELECT current_database()").Scan(&currentDatabase).Error; err != nil {
		log.WithError(err).Error("failed to detect current database")
		return err
	}

	var currentUser string
	if err := db.Raw("SELECT current_user").Scan(&currentUser).Error; err != nil {
		log.WithError(err).Error("failed to detect current user")
		return err
	}

	if currentDatabase != config.AppConfig.Postgres.Database {
		err := fmt.Errorf("connected to '%s' but expected '%s'", currentDatabase, config.AppConfig.Postgres.Database)
		log.WithError(err).Error("database mismatch")
		return err
	}

	log.WithFields(map[string]any{
		"database": currentDatabase,
		"user":     currentUser,
	}).Info("connected to database")

	debugDB := db.Session(&gorm.Session{Logger: gormlogger.Default.LogMode(gormlogger.Info)})

	var currentSchema string
	if err := debugDB.Raw("SELECT current_schema()").Scan(&currentSchema).Error; err != nil {
		log.WithError(err).Error("failed to detect current schema")
		return err
	}
	log.WithField("schema", currentSchema).Info("using schema")

	if err := debugDB.AutoMigrate(migrationModels...); err != nil {
		log.WithError(err).Error("failed to run auto migrate")
		return err
	}

	for _, m := range migrationModels {
		tableName, err := tableNameFromModel(db, m)
		if err != nil {
			log.WithError(err).Error("failed to resolve table name")
			return err
		}

		if !debugDB.Migrator().HasTable(tableName) {
			err := fmt.Errorf("table '%s' was not created for model %T", tableName, m)
			log.WithError(err).Error("migration verification failed")
			return err
		}

		log.WithField("table", tableName).Info("table verified")
	}

	log.Info("up migration completed")

	return nil
}

func Down(configPath *string, sslMode string) error {
	log := config.NewLogrusWithCategory("migration")
	log.Info("running down migration")

	if _, err := config.LoadConfig(configPath); err != nil {
		log.WithError(err).Error("failed to load config")
		return err
	}

	sslEnabled := sslMode == "enable"

	if config.AppConfig == nil || config.AppConfig.Postgres.Database == "" {
		err := fmt.Errorf("postgres database is not configured")
		log.WithError(err).Error("failed to resolve target database")
		return err
	}

	log.WithFields(map[string]any{
		"host":     config.AppConfig.Postgres.Host,
		"port":     config.AppConfig.Postgres.Port,
		"user":     config.AppConfig.Postgres.User,
		"database": config.AppConfig.Postgres.Database,
		"ssl":      sslEnabled,
	}).Info("target postgres connection")

	db, err := config.InitPostgresDBWithSSL(sslEnabled)
	if err != nil {
		log.WithError(err).Error("failed to initialize postgres")
		return err
	}
	defer config.ClosePostgresDB()

	var currentDatabase string
	if err := db.Raw("SELECT current_database()").Scan(&currentDatabase).Error; err != nil {
		log.WithError(err).Error("failed to detect current database")
		return err
	}

	if currentDatabase != config.AppConfig.Postgres.Database {
		err := fmt.Errorf("connected to '%s' but expected '%s'", currentDatabase, config.AppConfig.Postgres.Database)
		log.WithError(err).Error("database mismatch")
		return err
	}

	for i := len(migrationModels) - 1; i >= 0; i-- {
		tableName, err := tableNameFromModel(db, migrationModels[i])
		if err != nil {
			log.WithError(err).Error("failed to resolve table name")
			return err
		}

		if err := db.Migrator().DropTable(tableName); err != nil {
			log.WithError(err).Error("failed to drop table")
			return err
		}

		log.WithField("table", tableName).Info("table dropped")
	}

	log.Info("down migration completed")

	return nil
}

func tableNameFromModel(db *gorm.DB, m any) (string, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(m); err != nil {
		return "", err
	}

	if stmt.Schema == nil || stmt.Schema.Table == "" {
		return "", fmt.Errorf("unable to parse schema for model %T", m)
	}

	return stmt.Schema.Table, nil
}
