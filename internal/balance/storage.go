package balance

import "context"

type Repository interface {
	Create(ctx context.Context, balance *Balance) error
	FindOne(ctx context.Context, id uint64) (Balance, error)
	AddAmount(ctx context.Context, dto CreateUserBalanceDTO) error
	SubtractAmount(ctx context.Context, dto CreateUserBalanceDTO) error
}
