package reserve

import "context"

type Repository interface {
	Create(ctx context.Context, reserve *Reserve) error
	Delete(ctx context.Context, reserve *Reserve) error
	UpdateProfit(ctx context.Context, reserve *Reserve) error
}
