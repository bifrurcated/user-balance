package balance

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

func (s *Service) GetOne(ctx context.Context, id uint64) (user Balance, err error) {
	user, err = s.repository.FindOne(ctx, id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *Service) AddAmount(ctx context.Context, tum TransferUserMoney) error {
	err := s.repository.AddAmount(ctx, tum)
	if err != nil {
		return err
	}
	return err
}
