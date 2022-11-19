package user

import (
	"context"
	"github.com/bifrurcated/user-balance/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) *Service {
	return &Service{repository: repository, logger: logger}
}

func (s *Service) GetOne(ctx context.Context, id int64) (user User, err error) {
	user, err = s.repository.FindOne(ctx, id)
	if err != nil {
		return user, err
	}
	return user, nil
}
