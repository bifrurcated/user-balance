package history

import (
	"context"
	"github.com/bifrurcated/user-balance/pkg/api/sort"
)

type Repository interface {
	Create(ctx context.Context, history *History) error
	FindByUserID(ctx context.Context, userID uint64, options sort.Options) ([]History, error)
}
