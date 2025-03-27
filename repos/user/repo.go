package user

import (
	"time"
	"context"

	"github.com/google/uuid"
)

type Repo interface {
	Get(ctx context.Context, email string) (*User, error)
	Create(
		ctx context.Context,
		fullName string,
		email string,
		hashedPassword string,
	) error
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	FullName         string
	Email            string
	EmailConfirmedAt *time.Time
	HashedPassword   string
}
