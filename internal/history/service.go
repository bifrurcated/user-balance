package history

import (
	"context"
	"github.com/bifrurcated/user-balance/pkg/api/sort"
	"github.com/bifrurcated/user-balance/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) *Service {
	return &Service{repository: repository, logger: logger}
}

func (s *Service) UserTransactions(ctx context.Context, userID uint64, options sort.Options) ([]History, error) {
	histories, err := s.repository.FindByUserID(ctx, userID, options)
	if err != nil {
		return nil, err
	}
	return histories, nil
}
