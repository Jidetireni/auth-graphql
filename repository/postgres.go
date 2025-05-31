package repository

import (
	"auth-graphql/config"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func DatabaseInit(c *config.Config) (*sql.DB, error) {
	url := c.GetDatabaseURL() // get the database URL from the config
	if url == "" {
		return nil, fmt.Errorf("database URL is not set")
	}
	// Open a new database connection
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection %v", err)
	}

	// Set connection pool settings
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return db, nil

}
