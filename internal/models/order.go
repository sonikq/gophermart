package models

import "time"

type Order struct {
	Number     string    `json:"number"`
	Username   string    `json:"-"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
	Withdrawn  *float64  `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ListOrdersResponse []Order
