package database

import (
	"betterreads/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Book{})
	if err != nil {
		return err
	}
	return nil
}

func GetDbHandler() (*gorm.DB, error) {
	connectionString, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("Environment variable DATABASE_URL is not set, Must be a filename for sqlite or a connection string for postgres")
	}
	dbType, ok := os.LookupEnv("DB_TYPE")
	if !ok {
		panic("Environment variable DB_TYPE is not set, Must be 'sqlite' or 'postgres'")
	}

	var dialector gorm.Dialector

	switch {
	case dbType == "cosmos":
		dialector = postgres.Open(connectionString)
	case dbType == "sqlite":
		dialector = sqlite.Open(connectionString)
	default:
		panic("Unsupported database type")
	}

	dbHandler, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	err = migrate(dbHandler)
	if err != nil {
		return nil, err
	}

	return dbHandler, nil
}

func GetTestDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Couldn't create in-memory sqlite database for testing")
	}

	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		return nil
	}
	return db
}
