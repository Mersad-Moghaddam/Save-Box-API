package models

type Transaction struct {
	ID     string  `json:"id"`
	BoxID  string  `json:"boxId"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}
