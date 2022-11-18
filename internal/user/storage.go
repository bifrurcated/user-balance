package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindOne(ctx context.Context, id int64) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id int64) error
}
