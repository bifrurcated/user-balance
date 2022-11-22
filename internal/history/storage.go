package history

import "context"

type Repository interface {
	Create(ctx context.Context, history *History) error
}
