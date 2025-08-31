package dto

type CreateTableRequest struct {
	Seat   int    `json:"seat"`
	Status string `json:"status"` // available|occupied|closed (optional, default available)
}

type UpdateTableStatusRequest struct {
	Status string `json:"status" validate:"required"` // available|occupied|closed
}
