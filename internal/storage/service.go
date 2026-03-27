// Package storage provides a persistence layer using SQLite to maintain
// a history of all deployment attempts and their outcomes.
package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

// Deployment represents a single record of a deployment attempt in the database.
type Deployment struct {
	// ID is the unique auto-incrementing identifier for the record.
	ID int
	// Hash is the Git commit SHA detected for this deployment.
	Hash string
	// Status is the result of the operation (success, failed, or rollback).
	Status string
	// Output is the combined stdout/stderr from the remote execution.
	Output string
	// CreatedAt is the timestamp when the record was inserted.
	CreatedAt time.Time
}

// Service manages the connection pool and SQL operations for the metadata store.
type Service struct {
	db *sql.DB
}

// NewService initializes the SQLite database at the specified path,
// runs internal migrations, and returns a usable [Service] instance.
func NewService(dbPath string) (*Service, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	s := &Service{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}

	return s, nil
}

// migrate ensures the deployments table exists with the correct schema.
func (s *Service) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS deployments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hash TEXT NOT NULL,
		status TEXT NOT NULL,
		output TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := s.db.Exec(query)
	return err
}

// RecordDeploy inserts a new deployment event into the database history.
func (s *Service) RecordDeploy(hash, status, output string) error {
	query := `INSERT INTO deployments (hash, status, output) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, hash, status, output)
	return err
}

// GetHistory retrieves the last N deployment records from the database,
// ordered by their IDs in descending order.
func (s *Service) GetHistory(limit int) ([]Deployment, error) {
	query := `SELECT id, hash, status, output, created_at FROM deployments ORDER BY id DESC LIMIT ?`
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []Deployment
	for rows.Next() {
		var d Deployment
		if err := rows.Scan(&d.ID, &d.Hash, &d.Status, &d.Output, &d.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, d)
	}

	return history, nil
}

// Close terminates the underlying database connection.
func (s *Service) Close() error {
	return s.db.Close()
}
