package historydb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/history"
	"github.com/bifrurcated/user-balance/pkg/client/postgresql"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, history *history.History) error {
	q := `
		INSERT INTO history_operations 
		    (sender_user_id, user_id, service_id, amount, type) 
		VALUES 
		    ($1, $2, $3, $4, $5)
		RETURNING id, datetime
	`
	err := r.client.QueryRow(ctx, q,
		history.SenderUserID,
		history.UserID,
		history.ServiceID,
		history.Amount,
		history.Type).Scan(&history.ID, &history.Datetime)
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

func NewRepository(client postgresql.Client, logger *logging.Logger) history.Repository {
	return &repository{client: client, logger: logger}
}
