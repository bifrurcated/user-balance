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

func (s *Service) ReserveMoney(ctx context.Context, dto *CreateReserveDTO) (*Reserve, error) {
	err := s.balanceRepository.SubtractAmount(ctx, balance.CreateUserBalanceDTO{
		UserID: dto.UserID,
		Amount: dto.Cost,
	})
	if err != nil {
		return nil, err
	}
	reserve := NewReserve(dto)
	err = s.reserveRepository.Create(ctx, reserve)
	if err != nil {
		return nil, err
	}
	return reserve, err
}

func (s *Service) Delete(ctx context.Context, dto *CreateReserveDTO) (*Reserve, error) {
	reserve := NewReserve(dto)
	err := s.reserveRepository.Delete(ctx, reserve)
	if err != nil {
		return nil, err
	}
	return reserve, nil
}

func (s *Service) CancelReserve(ctx context.Context, dto *CancelReserveDTO) error {
	reserve := NewReserve(&CreateReserveDTO{
		UserID:    dto.UserID,
		ServiceID: dto.ServiceID,
		OrderID:   dto.OrderID,
	})
	err := s.reserveRepository.Delete(ctx, reserve)
	if err != nil {
		return err
	}
	err = s.balanceRepository.AddAmount(ctx, balance.CreateUserBalanceDTO{
		UserID: reserve.UserID,
		Amount: reserve.Cost,
	})
	if err != nil {
		return err
	}
	return nil
}
