package models

import "time"

type Product struct {
	ID             int       `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Code           string    `db:"code" json:"code"`
	Description    string    `db:"description" json:"description"`
	InterestRate   float64   `db:"interest_rate" json:"interest_rate"`
	MinAmount      float64   `db:"min_amount" json:"min_amount"`
	MaxAmount      float64   `db:"max_amount" json:"max_amount"`
	TermMinMonths  int       `db:"term_min_months" json:"term_min_months"`
	TermMaxMonths  int       `db:"term_max_months" json:"term_max_months"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	TypeID         int       `db:"type_id" json:"type_id"`
	ModuleID       int       `db:"module_id" json:"module_id"`
	StateID        int       `db:"state_id" json:"state_id"`
}
