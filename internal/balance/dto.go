package balance

type TransferUserMoney struct {
	UserID uint64  `json:"user_id"`
	Amount float32 `json:"amount"`
}

type UserBalance struct {
	UserID uint64 `json:"user_id"`
}
