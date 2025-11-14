package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// PostgresDB представляет подключение к PostgreSQL
type PostgresDB struct {
	db *sql.DB
}

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB(connStr string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{db: db}, nil
}

// Close закрывает подключение к БД
func (p *PostgresDB) Close() error {
	return p.db.Close()
}

// DB возвращает *sql.DB для использования в репозиториях
func (p *PostgresDB) DB() *sql.DB {
	return p.db
}
