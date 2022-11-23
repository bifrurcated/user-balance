package history

type CreateHistoryDTO struct {
	SenderUserID *uint64 `json:"sender_user_id"`
	UserID       uint64  `json:"user_id"`
	ServiceID    *uint64 `json:"service_id"`
	Amount       float32 `json:"amount"`
	Type         string  `json:"type"`
}

type UserHistoryDTO struct {
	UserID uint64 `json:"user_id"`
}

type OptionsDTO struct {
	Limit uint64
	Value any
	Field string
	Order string
}
