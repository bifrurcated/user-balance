package balancedb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/apperror"
	"github.com/bifrurcated/user-balance/internal/balance"
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

func (r *repository) Create(ctx context.Context, balance *balance.Balance) error {
	q := `
			INSERT INTO balance 
			    (user_id, amount) 
			VALUES 
			    ($1,$2)
			RETURNING id
		`
	r.logger.Tracef("SQL Query: %s", repeatable.FormatQuery(q))
	err := r.client.QueryRow(ctx, q, balance.UserID, balance.Amount).Scan(&balance.ID)
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

func (r *repository) FindOne(ctx context.Context, id uint64) (u balance.Balance, err error) {
	q := `
		SELECT id, user_id, amount FROM balance WHERE user_id=$1
	`
	r.logger.Tracef("SQL Query: %s", repeatable.FormatQuery(q))
	err = r.client.QueryRow(ctx, q, id).Scan(&u.ID, &u.UserID, &u.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, apperror.ErrNotFound
		}
		return u, err
	}

	return u, nil
}

func (r *repository) AddAmount(ctx context.Context, tum balance.TransferUserMoney) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	q := `
		UPDATE balance
		SET amount = amount + $2
		WHERE user_id = $1
	`
	r.logger.Tracef("SQL Query: %s", repeatable.FormatQuery(q))
	res, err := r.client.Exec(ctx, q, tum.UserID, tum.Amount)
	if err != nil {
		return err
	}
	r.logger.Debug(res)
	r.logger.Debugf("RowsAffected: %d", res.RowsAffected())
	if res.RowsAffected() == 0 {
		err = r.Create(ctx, &balance.Balance{
			UserID: tum.UserID,
			Amount: tum.Amount,
		})
		if err != nil {
			return err
		}
	}
	tx.Commit(ctx)
	return nil
}

func (r *repository) SubtractAmount(ctx context.Context, tum balance.TransferUserMoney) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	one, err := r.FindOne(ctx, tum.UserID)
	if err != nil {
		return err
	}
	if (one.Amount - tum.Amount) < 0 {
		return apperror.NewAppError(nil,
			fmt.Sprintf("not enough money (%f) on the user (%d) balance, required amount: %f", one.Amount, tum.UserID, tum.Amount),
			"", "US-000005")
	}
	q := `
		UPDATE balance
		SET amount = amount - $2
		WHERE user_id = $1
	`
	res, err := r.client.Exec(ctx, q, tum.UserID, tum.Amount)
	if err != nil {
		return err
	}
	r.logger.Debug(res)
	tx.Commit(ctx)
	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) balance.Repository {
	return &repository{client: client, logger: logger}
}
