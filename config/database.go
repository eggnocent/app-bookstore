package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Postgres struct {
	DB *sqlx.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Psql.User,
		cfg.Psql.Password,
		cfg.Psql.Host,
		cfg.Psql.Port,
		cfg.Psql.DBName,
	)

	db, err := sqlx.Connect("postgres", dbConnString)
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-1] Failed to connect to database")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-2] Database connection test failed")
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Psql.DBMaxOpen)
	db.SetMaxIdleConns(cfg.Psql.DBMaxIdle)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Info().Msg("[ConnectionPostgres] Database connection established successfully")
	return &Postgres{DB: db}, nil
}
