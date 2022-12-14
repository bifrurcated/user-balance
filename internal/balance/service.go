package balance

import (
	"context"
	"github.com/bifrurcated/user-balance/internal/history"
	"github.com/bifrurcated/user-balance/pkg/logging"
)

type Service struct {
	balanceRepository Repository
	historyRepository history.Repository
	logger            *logging.Logger
}

func NewService(repository Repository, historyRepository history.Repository, logger *logging.Logger) *Service {
	return &Service{balanceRepository: repository, historyRepository: historyRepository, logger: logger}
}

func (s *Service) GetOne(ctx context.Context, id uint64) (user Balance, err error) {
	user, err = s.balanceRepository.FindOne(ctx, id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *Service) AddAmount(ctx context.Context, dto CreateUserBalanceDTO) error {
	err := s.balanceRepository.AddAmount(ctx, dto)
	if err != nil {
		return err
	}
	err = s.historyRepository.Create(ctx, &history.History{
		UserID: dto.UserID,
		Amount: dto.Amount,
		Type:   "зачисление",
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) TransferUserMoney(ctx context.Context, dto TransferUserMoneyDTO) error {
	err := s.balanceRepository.SubtractAmount(ctx, CreateUserBalanceDTO{
		UserID: dto.SenderUserID,
		Amount: dto.Amount,
	})
	if err != nil {
		return err
	}
	err = s.balanceRepository.AddAmount(ctx, CreateUserBalanceDTO{
		UserID: dto.ReceiverUserID,
		Amount: dto.Amount,
	})
	if err != nil {
		return err
	}
	err = s.historyRepository.Create(ctx, &history.History{
		SenderUserID: &dto.SenderUserID,
		UserID:       dto.ReceiverUserID,
		Amount:       dto.Amount,
		Type:         "перевод",
	})
	if err != nil {
		return err
	}

	return nil
}
