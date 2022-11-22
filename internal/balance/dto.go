package balance

type CreateUserBalanceDTO struct {
	UserID uint64  `json:"user_id"`
	Amount float32 `json:"amount"`
}

type TransferUserMoneyDTO struct {
	SenderUserID   uint64  `json:"sender_user_id"`
	ReceiverUserID uint64  `json:"receiver_user_id"`
	Amount         float32 `json:"amount"`
}

type UserBalance struct {
	UserID uint64 `json:"user_id"`
}
