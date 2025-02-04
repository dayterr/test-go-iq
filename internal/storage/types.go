package storage

const (
	Arrival  = "arrival"
	Transfer = "transfer"
)

type Storage struct {
	DSN string
}

type User struct {
	ID      int     `json:"id"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	ID         int
	SenderID   int     `json:"sender_id"`
	Amount     float64 `json:"amount"`
	RecieverID int     `json:"reciever_id"`
}

type Operation struct {
	ID            int
	UserID        int
	OperationType string
	Amount        float64
}
