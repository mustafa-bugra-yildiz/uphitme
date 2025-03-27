package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -typed -package user -destination mock.go . Repo
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
