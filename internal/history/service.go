package history

import (
	"github.com/bifrurcated/user-balance/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewService(repository Repository, logger *logging.Logger) *Service {
	return &Service{repository: repository, logger: logger}
}
