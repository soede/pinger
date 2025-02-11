package postgres

import (
	"docker/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s sslrootcert=%s sslcert=%s sslkey=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Name,
		c.DB.SSL,
		c.DB.Password,
		c.DB.SSLRootCert,
		c.DB.SSLCert,
		c.DB.SSLKey,
	)

	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
