package models

type WithdrawRequest struct {
	Order    string  `json:"order"`
	Sum      float64 `json:"sum"`
	Username string  `json:"-"`
}

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
