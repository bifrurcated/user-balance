package reserve

import "github.com/bifrurcated/user-balance/internal/user"

type Reserve struct {
	ID        int64     `json:"id"`
	User      user.User `json:"user"`
	ServiceID int64     `json:"service_id"`
	OrderID   int64     `json:"order_id"`
	Amount    float32   `json:"amount"`
}
