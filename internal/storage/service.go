package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

// Deployment representa um registro de deploy no banco de dados.
type Deployment struct {
	ID        int
	Hash      string
	Status    string
	Output    string
	CreatedAt time.Time
}

// Service gerencia o acesso ao banco de dados SQLite.
type Service struct {
	db *sql.DB
}

// NewService cria e inicializa o servico de storage.
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

// migrate cria as tabelas necessarias se nao existirem.
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

// RecordDeploy registra uma nova tentativa de deploy.
func (s *Service) RecordDeploy(hash, status, output string) error {
	query := `INSERT INTO deployments (hash, status, output) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, hash, status, output)
	return err
}

// GetHistory retorna os ultimos N deploys.
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

// Close fecha a conexao com o banco de dados.
func (s *Service) Close() error {
	return s.db.Close()
}
