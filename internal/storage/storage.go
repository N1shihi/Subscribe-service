package storage

import (
	"Subscribe-service/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func New(cfg *config.Config) (*sql.DB, error) {
	conStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SslMode,
	)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	const (
		maxAttempts = 10
		delay       = 2 * time.Second
	)

	var pingErr error
	for i := 1; i <= maxAttempts; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			log.Printf("db ping succeeded on attempt %d", i)
			return db, nil
		}

		log.Printf("db ping failed (attempt %d/%d): %v", i, maxAttempts, pingErr)
		time.Sleep(delay)
	}

	_ = db.Close()
	return nil, fmt.Errorf("ping db: %w", pingErr)
}
