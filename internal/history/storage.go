package history

import "context"

type Repository interface {
	Create(ctx context.Context, history *History) error
	FindByUserID(ctx context.Context, userID uint64) ([]History, error)
}
