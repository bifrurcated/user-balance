package balance

import "context"

type Repository interface {
	Create(ctx context.Context, balance *Balance) error
	FindOne(ctx context.Context, id int64) (Balance, error)
	AddAmount(ctx context.Context, tum TransferUserMoney) error
}
