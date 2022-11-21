package balance

type Balance struct {
	ID     uint64  `json:"id"`
	UserID uint64  `json:"user_id"`
	Amount float32 `json:"amount"`
}
