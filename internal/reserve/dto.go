package reserve

type CreateReserveDTO struct {
	UserID    uint64  `json:"user_id"`
	ServiceID uint64  `json:"service_id"`
	OrderID   uint64  `json:"order_id"`
	Cost      float32 `json:"cost"`
	IsProfit  bool    `json:"is_profit"`
}

type ProfitReserveDTO struct {
	ID        uint64  `json:"id"`
	UserID    uint64  `json:"user_id"`
	ServiceID uint64  `json:"service_id"`
	OrderID   uint64  `json:"order_id"`
	Cost      float32 `json:"cost"`
	Amount    float32 `json:"amount"`
}

type CancelReserveDTO struct {
	UserID    uint64 `json:"user_id"`
	ServiceID uint64 `json:"service_id"`
	OrderID   uint64 `json:"order_id"`
}

type CreateReserveMoneyDTO struct {
	UserID    uint64  `json:"user_id"`
	ServiceID uint64  `json:"service_id"`
	OrderID   uint64  `json:"order_id"`
	Cost      float32 `json:"cost"`
}
