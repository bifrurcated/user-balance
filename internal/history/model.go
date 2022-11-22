package history

import "time"

type History struct {
	ID           uint64    `json:"id"`
	SenderUserID *uint64   `json:"sender_user_id"`
	UserID       uint64    `json:"user_id"`
	ServiceID    *uint64   `json:"service_id"`
	Amount       float32   `json:"amount"`
	Type         string    `json:"type"`
	Datetime     time.Time `json:"datetime"`
}

func NewHistory(dto *CreateHistoryDTO) *History {
	return &History{
		SenderUserID: dto.SenderUserID,
		UserID:       dto.UserID,
		ServiceID:    dto.ServiceID,
		Amount:       dto.Amount,
		Type:         dto.Type,
	}
}
