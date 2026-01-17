package dto

// ReservationMessage is the message format for reservation queue
type ReservationMessage struct {
	ReservationID  string `json:"reservation_id"`
	UserID         string `json:"user_id"`
	PointAmount    int    `json:"point_amount"`
	IdempotencyKey string `json:"idempotency_key"`
}
