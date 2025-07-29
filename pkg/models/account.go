package models

type Account struct {
	ID            int     `db:"id" json:"id"`
	AccountNumber string  `db:"account_number" json:"account_number"`
	Balance       float64 `db:"account_balance" json:"account_balance"`
	Cards         []Card  `json:"cards,omitempty"`
}
