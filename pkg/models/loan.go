package models
import "time"
type Loan struct {
	ID              int       `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Amount          float64   `db:"amount" json:"amount"`
	RemainingAmount float64   `db:"remaining_amount" json:"remaining_amount"`
	InterestRate    float64   `db:"interest_rate" json:"interest_rate"`
	TermMonths      int       `db:"term_months" json:"term_months"`
	StartDate       time.Time `db:"start_date" json:"start_date"`
	EndDate         time.Time `db:"end_date" json:"end_date"`
	StateCode       string    `db:"state_code" json:"state_code"`
	ProductName     string    `db:"product_name" json:"product_name"`
	PaymentSchedule []Payment `json:"payment_schedule"`
}