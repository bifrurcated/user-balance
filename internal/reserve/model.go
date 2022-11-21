package reserve

type Reserve struct {
	ID        uint64  `json:"id"`
	UserID    uint64  `json:"user_id"`
	ServiceID uint64  `json:"service_id"`
	OrderID   uint64  `json:"order_id"`
	Cost      float32 `json:"cost"`
}

func NewReserve(dto *CreateReserveDTO) *Reserve {
	return &Reserve{
		UserID:    dto.UserID,
		ServiceID: dto.ServiceID,
		OrderID:   dto.OrderID,
		Cost:      dto.Cost,
	}
}
