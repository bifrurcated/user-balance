package db

import (
	"context"
	"errors"
	"github.com/bifrurcated/user-balance/internal/apperror"
	"github.com/bifrurcated/user-balance/internal/user"
	"github.com/bifrurcated/user-balance/pkg/client/postgresql"
	"github.com/bifrurcated/user-balance/pkg/logging"
	repeatable "github.com/bifrurcated/user-balance/pkg/utils"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, user *user.User) error {
	panic("")
}

func (r *repository) FindOne(ctx context.Context, id int64) (u user.User, err error) {
	q := `
		SELECT id, name, amount FROM usr WHERE id=$1
	`
	r.logger.Tracef("SQL Query: %s", repeatable.FormatQuery(q))
	err = r.client.QueryRow(ctx, q, id).Scan(&u.ID, &u.Name, &u.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, apperror.ErrNotFound
		}
		return u, err
	}

	return u, nil
}

func (r *repository) FindAll(ctx context.Context) ([]user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) user.Repository {
	return &repository{client: client, logger: logger}
}
