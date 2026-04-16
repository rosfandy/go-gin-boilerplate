package migration

import (
	"fmt"
	"owner-api-proxy/internal/config"
	"owner-api-proxy/internal/database/model"
	"strings"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var migrationModels = []any{
	&model.Users{},
	&model.Owner{},
	&model.OwnerPremium{},
}

func Up(configPath *string, sslMode string) error {
	return UpWithDB(configPath, sslMode, "postgres")
}

func UpWithDB(configPath *string, sslMode string, dbType string) error {
	log := config.NewLogrusWithCategory("migration")
	log.Info("running up migration")

	if _, err := config.LoadConfig(configPath); err != nil {
		log.WithError(err).Error("failed to load config")
		return err
	}

	normalizedDBType := strings.ToLower(strings.TrimSpace(dbType))
	dbConfig, err := getDBConfig(normalizedDBType)
	if err != nil {
		log.WithError(err).Error("failed to resolve target database")
		return err
	}

	if sslMode != "" {
		dbConfig.Secure = sslMode
	}

	if dbConfig.Database == "" {
		err := fmt.Errorf("%s database is not configured", normalizedDBType)
		log.WithError(err).Error("failed to resolve target database")
		return err
	}

	log.WithFields(map[string]any{
		"db_type":  normalizedDBType,
		"host":     dbConfig.Host,
		"port":     dbConfig.Port,
		"user":     dbConfig.User,
		"database": dbConfig.Database,
		"ssl":      dbConfig.Secure,
	}).Info("target database connection")

	dbClient := config.NewClientSql(normalizedDBType)
	dbErr := dbClient.Open()
	if dbErr != nil {
		log.WithError(dbErr).Error("failed to initialize database")
		return dbErr
	}
	defer dbClient.Close()

	db := dbClient.Client()
	if db == nil {
		err := fmt.Errorf("database client is not initialized")
		log.WithError(err).Error("failed to initialize database")
		return err
	}

	currentDatabaseQuery := "SELECT current_database()"
	if normalizedDBType == "mysql" {
		currentDatabaseQuery = "SELECT DATABASE()"
	}

	var currentDatabase string
	if err := db.Raw(currentDatabaseQuery).Scan(&currentDatabase).Error; err != nil {
		log.WithError(err).Error("failed to detect current database")
		return err
	}

	if currentDatabase != dbConfig.Database {
		err := fmt.Errorf("connected to '%s' but expected '%s'", currentDatabase, dbConfig.Database)
		log.WithError(err).Error("database mismatch")
		return err
	}

	var currentUser string
	if err := db.Raw("SELECT current_user").Scan(&currentUser).Error; err != nil {
		log.WithError(err).Error("failed to detect current user")
		return err
	}

	log.WithFields(map[string]any{
		"db_type":  normalizedDBType,
		"database": currentDatabase,
		"user":     currentUser,
	}).Info("connected to database")

	debugDB := db.Session(&gorm.Session{Logger: gormlogger.Default.LogMode(gormlogger.Info)})

	if normalizedDBType == "postgres" {
		var currentSchema string
		if err := debugDB.Raw("SELECT current_schema()").Scan(&currentSchema).Error; err != nil {
			log.WithError(err).Error("failed to detect current schema")
			return err
		}
		log.WithField("schema", currentSchema).Info("using schema")
	}

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
	return DownWithDB(configPath, sslMode, "postgres")
}

func DownWithDB(configPath *string, sslMode string, dbType string) error {
	log := config.NewLogrusWithCategory("migration")
	log.Info("running down migration")

	if _, err := config.LoadConfig(configPath); err != nil {
		log.WithError(err).Error("failed to load config")
		return err
	}

	normalizedDBType := strings.ToLower(strings.TrimSpace(dbType))
	dbConfig, err := getDBConfig(normalizedDBType)
	if err != nil {
		log.WithError(err).Error("failed to resolve target database")
		return err
	}

	if sslMode != "" {
		dbConfig.Secure = sslMode
	}

	if dbConfig.Database == "" {
		err := fmt.Errorf("%s database is not configured", normalizedDBType)
		log.WithError(err).Error("failed to resolve target database")
		return err
	}

	log.WithFields(map[string]any{
		"db_type":  normalizedDBType,
		"host":     dbConfig.Host,
		"port":     dbConfig.Port,
		"user":     dbConfig.User,
		"database": dbConfig.Database,
		"ssl":      dbConfig.Secure,
	}).Info("target database connection")

	dbClient := config.NewClientSql(normalizedDBType)
	dbErr := dbClient.Open()
	if dbErr != nil {
		log.WithError(dbErr).Error("failed to initialize database")
		return dbErr
	}
	defer dbClient.Close()

	db := dbClient.Client()
	if db == nil {
		err := fmt.Errorf("database client is not initialized")
		log.WithError(err).Error("failed to initialize database")
		return err
	}

	currentDatabaseQuery := "SELECT current_database()"
	if normalizedDBType == "mysql" {
		currentDatabaseQuery = "SELECT DATABASE()"
	}

	var currentDatabase string
	if err := db.Raw(currentDatabaseQuery).Scan(&currentDatabase).Error; err != nil {
		log.WithError(err).Error("failed to detect current database")
		return err
	}

	if currentDatabase != dbConfig.Database {
		err := fmt.Errorf("connected to '%s' but expected '%s'", currentDatabase, dbConfig.Database)
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

func getDBConfig(dbType string) (*config.DbConfig, error) {
	if config.AppConfig == nil {
		return nil, fmt.Errorf("app config is not loaded")
	}

	switch dbType {
	case "postgres":
		return &config.AppConfig.Postgres, nil
	case "mysql":
		return &config.AppConfig.MySQL, nil
	default:
		return nil, fmt.Errorf("unsupported db type: %s", dbType)
	}
}
