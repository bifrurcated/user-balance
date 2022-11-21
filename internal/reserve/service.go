package reserve

import (
	"context"
	"github.com/bifrurcated/user-balance/internal/balance"
	"github.com/bifrurcated/user-balance/pkg/logging"
)

type Service struct {
	reserveRepository Repository
	balanceRepository balance.Repository
	logger            *logging.Logger
}

func NewService(reserveRepository Repository, balanceRepository balance.Repository, logger *logging.Logger) *Service {
	return &Service{reserveRepository: reserveRepository, balanceRepository: balanceRepository, logger: logger}
}

func (s *Service) ReserveMoney(ctx context.Context, dto CreateReserveDTO) (*Reserve, error) {
	err := s.balanceRepository.SubtractAmount(ctx, balance.TransferUserMoney{
		UserID: dto.UserID,
		Amount: dto.Amount,
	})
	if err != nil {
		return nil, err
	}
	reserve := NewReserve(&dto)
	err = s.reserveRepository.Create(ctx, reserve)
	if err != nil {
		return nil, err
	}
	return reserve, err
}
