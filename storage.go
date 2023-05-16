package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStory() (*PostgresStore, error) {
	connStr := "user=admin dbname=gobank password=admin sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err

	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			number SERIAL,
			balance BIGINT,
			create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func (s *PostgresStore) CreateAccount(a *Account) error {
	query := `
		INSERT INTO accounts (first_name, last_name, number, balance)
		VALUES ($1, $2, $3, $4)
		RETURNING id, create_at;
	`
	resp, err := s.db.Query(
		query,
		a.FirstName,
		a.LastName,
		a.Number,
		a.Balance,
	)
	if err != nil {
		return err
	}

	defer resp.Close()

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`
		SELECT id, first_name, last_name, number, balance, create_at
		FROM accounts;
	`)
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0)
	for rows.Next() {
		var a Account
		if err := rows.Scan(
			&a.ID,
			&a.FirstName,
			&a.LastName,
			&a.Number,
			&a.Balance,
			&a.CreatedAt,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, &a)
	}

	return accounts, nil
}
