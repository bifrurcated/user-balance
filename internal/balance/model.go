package balance

type Balance struct {
	ID     int64   `json:"id"`
	UserID int64   `json:"user_id"`
	Amount float32 `json:"amount"`
}
