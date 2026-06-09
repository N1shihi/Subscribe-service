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
		return nil, fmt.Errorf("db open error: %w", err)
	}

	const (
		maxAttempts = 10
		delay       = 2 * time.Second
	)

	for i := 1; i <= maxAttempts; i++ {
		if err := db.Ping(); err == nil {
			log.Printf("db connected after %d attempts", i)
			return db, nil
		} else {
			log.Printf("db connection attempt %d failed", i)
		}

		time.Sleep(delay)
	}

	_ = db.Close()
	return nil, fmt.Errorf("database is unreachable after %d attempts", maxAttempts)
}
