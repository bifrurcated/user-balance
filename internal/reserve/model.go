package reserve

type Reserve struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	ServiceID int64   `json:"service_id"`
	OrderID   int64   `json:"order_id"`
	Amount    float32 `json:"amount"`
}
