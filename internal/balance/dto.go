package balance

type TransferUserMoney struct {
	UserID int64   `json:"user_id"`
	Amount float32 `json:"amount"`
}

type UserBalance struct {
	UserID int64 `json:"user_id"`
}
