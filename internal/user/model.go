package user

type User struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
}
