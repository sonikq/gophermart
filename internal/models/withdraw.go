package models

import "time"

type WithdrawRequest struct {
	Order    string  `json:"order"`
	Sum      float64 `json:"sum"`
	Username string  `json:"-"`
}

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
