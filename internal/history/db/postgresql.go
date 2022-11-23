package historydb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/apperror"
	"github.com/bifrurcated/user-balance/internal/history"
	"github.com/bifrurcated/user-balance/pkg/api/sort"
	"github.com/bifrurcated/user-balance/pkg/client/postgresql"
	"github.com/bifrurcated/user-balance/pkg/logging"
	repeatable "github.com/bifrurcated/user-balance/pkg/utils"
	"github.com/jackc/pgx/v5"
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
	r.logger.Tracef("SQL Query: %s", repeatable.FormatQuery(q))
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

func (r *repository) FindByUserID(ctx context.Context, userID uint64, options history.OptionsDTO) ([]history.History, error) {
	q := `
		SELECT id, sender_user_id, user_id, service_id, amount, type, datetime 
		FROM history_operations 
		WHERE (user_id = $1 OR sender_user_id = $1)
	`
	if options.Field != "" && options.Order != "" {
		//подставлять параметры сортировки в поле проверки
		//в зависимости от Order ставить больше или меньше
		sign := ">"
		if options.Order == sort.DESC {
			sign = "<"
		}
		if options.Value != nil {
			q = fmt.Sprintf("%s AND %s %s $3", q, options.Field, sign)
		}
		q = fmt.Sprintf("%s ORDER BY %s %s", q, options.Field, options.Order)
	}
	q = q + `
		LIMIT $2 
	`
	r.logger.Tracef("SQL Query: %s", repeatable.FormatQuery(q))
	var rows pgx.Rows
	var err error
	if options.Value != nil {
		rows, err = r.client.Query(ctx, q, userID, options.Limit, options.Value)
	} else {
		rows, err = r.client.Query(ctx, q, userID, options.Limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0
	historyArr := make([]history.History, 0)
	for rows.Next() {
		var h history.History
		err = rows.Scan(&h.ID, &h.SenderUserID, &h.UserID, &h.ServiceID, &h.Amount, &h.Type, &h.Datetime)
		if err != nil {
			return nil, err
		}
		historyArr = append(historyArr, h)
		count++
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, apperror.ErrNotFound
	}
	return historyArr, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) history.Repository {
	return &repository{client: client, logger: logger}
}
