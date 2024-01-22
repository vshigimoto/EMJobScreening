package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go-jun/internal/person/config"
)

type Config config.DbNode

func (c Config) dsn() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		c.Host,
		c.Port,
		c.DbName,
		c.User,
		c.Password,
	)
}

func New(cfg config.DbNode) (*sqlx.DB, error) {
	conf := Config(cfg)

	dbConn, err := sqlx.Connect("postgres", conf.dsn())
	if err != nil {
		return nil, err
	}
	if err := dbConn.Ping(); err != nil {
		return nil, err
	}
	return dbConn, nil
}
