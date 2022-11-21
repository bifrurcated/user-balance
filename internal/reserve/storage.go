package reserve

import "context"

type Repository interface {
	Create(ctx context.Context, reserve *Reserve) error
	FindOne(ctx context.Context, id uint64) (Reserve, error)
	Delete(ctx context.Context, reserve *Reserve) error
}
