package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return &postgres{db: db}
}

func (p postgres) Get(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT
			id, created_at, updated_at,
			full_name, email, email_confirmed_at, hashed_password
		FROM users
		WHERE
			id = $1
	`

	var user User
	err := p.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.FullName, &user.Email, &user.EmailConfirmedAt, &user.HashedPassword,
	)
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user.Repo.Get: %w", err)
	}

	return &user, nil
}

func (p postgres) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT
			id, created_at, updated_at,
			full_name, email, email_confirmed_at, hashed_password
		FROM users
		WHERE
			email = $1
	`

	var user User
	err := p.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
		&user.FullName, &user.Email, &user.EmailConfirmedAt, &user.HashedPassword,
	)
	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user.Repo.GetByEmail: %w", err)
	}

	return &user, nil
}

func (p postgres) Create(
	ctx context.Context,
	fullName string,
	email string,
	hashedPassword string,
) error {
	id := uuid.New()

	query := `
		INSERT INTO users
			(id, full_name, email, hashed_password)
		VALUES
			($1, $2, $3, $4)
	`

	_, err := p.db.Exec(ctx, query, id, fullName, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("user.Repo.Create: %w", err)
	}

	return nil
}
