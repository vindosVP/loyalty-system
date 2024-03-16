package models

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"strconv"
	"time"
)

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Status     string    `json:"status"`
	Sum        float64   `json:"sum"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (o *Order) Validate() error {
	return goluhn.Validate(strconv.Itoa(o.ID))
}
