package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/northwindman/testREST-autentification/internal/domain/models"
	"github.com/northwindman/testREST-autentification/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// Use this constant for initial place of the error(stack-trace)
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users
	(
		uid BIGSERIAL PRIMARY KEY,
		ip TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		pass_hash BYTEA NOT NULL,
		secret TEXT NOT NULL,
		refresh_token BYTEA NOT NULL
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_email ON users(email)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_secret ON users(secret);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser create new user in DB
func (s *Storage) SaveUser(ip string, email string, passHash []byte, secret string, refreshToken []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	query := `
		INSERT INTO users(ip, email, pass_hash, secret, refresh_token)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uid;
	`

	var uid int64
	err := s.db.QueryRow(query, ip, email, passHash, secret, refreshToken).Scan(&uid)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, storage.ErrAlreadyExist
			}
			fmt.Println(pqErr.Code)
		}
		fmt.Println(err.Error())

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}

// GetUser returns the user's model for the operation by email and secret
func (s *Storage) GetUser(email string) (models.User, error) {
	const op = "storage.postgres.GetUser"

	query := `
		SELECT uid, ip, email, pass_hash, secret, refresh_token
		FROM users
		WHERE email = $1;
	`

	var user models.User
	err := s.db.QueryRow(query, email).Scan(&user.UID, &user.IP, &user.Email, &user.PassHash, &user.Secret, &user.RefreshToken)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return models.User{}, storage.ErrNotFound
			}
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// UpdateUser updates the user's data , namely the refresh token and secret
func (s *Storage) UpdateUser(email string, ip string, secret string, refreshToken []byte) (int64, error) {
	const op = "storage.postgres.UpdateUser"

	// Вообще, по поводу обновления IP я не уверен, но всё зависит от логики приложения
	query := `
		UPDATE users
		SET
			ip = $1,
			secret = $2,
			refresh_token = $3
		WHERE
			email = $4
		RETURNING uid;
	`

	var uid int64
	err := s.db.QueryRow(query, ip, secret, refreshToken, email).Scan(&uid)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, storage.ErrNotFound
			}
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}
