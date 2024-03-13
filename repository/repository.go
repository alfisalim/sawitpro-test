// This file contains the repository implementation layer.
package repository

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

type NewRepositoryOptions struct {
	Dsn string
}

func NewRepository(opts NewRepositoryOptions) *Repository {
	sqlDB, err := sql.Open("pgx", opts.Dsn)
	if err != nil {
		panic("failed to connect database")
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlCon, _ := gormDB.DB()
	sqlCon.SetMaxIdleConns(10)
	sqlCon.SetMaxOpenConns(100)
	sqlCon.SetConnMaxLifetime(0)

	return &Repository{gormDB}
}
