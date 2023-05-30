package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
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
			password VARCHAR(255),
			create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func (s *PostgresStore) CreateAccount(a *Account) (*Account, error) {
	query := `
		INSERT INTO accounts (first_name, last_name, number, balance, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, create_at;
	`
	resp, err := s.db.Query(
		query,
		a.FirstName,
		a.LastName,
		a.Number,
		a.Balance,
		a.EncryptedPassword,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Close()

	if resp.Next() {
		err := resp.Scan(
			&a.ID,
			&a.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		return a, nil
	}

	return nil, fmt.Errorf("no account found")
}

func (s *PostgresStore) DeleteAccount(id int) error {
	res, err := s.db.Query(`
		DELETE FROM accounts
		WHERE id = $1;
		`, id)

	if err != nil {
		return err
	}

	defer res.Close()

	return nil
}

func (s *PostgresStore) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {

	rows, err := s.db.Query(`
		SELECT id, first_name, last_name, number, balance, create_at
		FROM accounts
		WHERE id = $1;
		`, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query(`
		SELECT id, first_name, last_name, number, balance, create_at, password
		FROM accounts
		WHERE number = $1;
		`, number)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		return scanIntoAccountWithPassword(rows)
	}

	return nil, fmt.Errorf("account with number %d not found", number)
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
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	var account Account
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func scanIntoAccountWithPassword(rows *sql.Rows) (*Account, error) {
	var account Account
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
		&account.EncryptedPassword,
	)

	if err != nil {
		return nil, err
	}

	return &account, nil
}
