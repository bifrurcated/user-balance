package history

type CreateHistoryDTO struct {
	SenderUserID *uint64 `json:"sender_user_id"`
	UserID       uint64  `json:"user_id"`
	ServiceID    *uint64 `json:"service_id"`
	Amount       float32 `json:"amount"`
	Type         string  `json:"type"`
}
