package reservedb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/reserve"
	"github.com/bifrurcated/user-balance/pkg/client/postgresql"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, reserve *reserve.Reserve) error {
	q := `
		INSERT INTO reserve 
		    (user_id, service_id, order_id, amount) 
		VALUES 
		    ($1,$2,$3,$4)
		RETURNING id
	`
	err := r.client.QueryRow(ctx, q, reserve.UserID, reserve.ServiceID, reserve.OrderID, reserve.Cost).Scan(&reserve.ID)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgError.Message,
				pgError.Detail,
				pgError.Where,
				pgError.Code,
				pgError.SQLState())
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindOne(ctx context.Context, id uint64) (reserve.Reserve, error) {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) reserve.Repository {
	return &repository{client: client, logger: logger}
}
