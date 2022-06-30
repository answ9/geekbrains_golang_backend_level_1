package userrepo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lesson9/reguser/internal/entities/userentity"
)

// нужен только тут
type UserStore interface {
	Create(ctx context.Context, u userentity.User) (*uuid.UUID, error)
	Read(ctx context.Context, uid uuid.UUID) (*userentity.User, error)
	Delete(ctx context.Context, uid uuid.UUID) error
	SearchUsers(ctx context.Context, s string) (chan userentity.User, error)
}

type Users struct {
	ustore UserStore
}

func NewUsers(ustore UserStore) *Users {
	return &Users{
		ustore: ustore,
	}
}

func (us *Users) Create(ctx context.Context, u userentity.User) (*userentity.User, error) {
	u.ID = uuid.New()
	id, err := us.ustore.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}
	u.ID = *id
	return &u, nil
}

func (us *Users) Read(ctx context.Context, uid uuid.UUID) (*userentity.User, error) {
	u, err := us.ustore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("read user error: %w", err)
	}
	return u, nil
}

func (us *Users) Delete(ctx context.Context, uid uuid.UUID) (*userentity.User, error) {
	u, err := us.ustore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("search user error: %w", err)
	}
	return u, us.ustore.Delete(ctx, uid)
}

func (us *Users) SearchUsers(ctx context.Context, s string) (chan userentity.User, error) {
	chin, err := us.ustore.SearchUsers(ctx, s)
	if err != nil {
		return nil, err
	}
	chout := make(chan userentity.User, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case u, ok := <-chin:
				if !ok {
					return
				}
				u.Permissions = 0755
				chout <- u
			}
		}
	}()
	return chout, nil
}
