package models
import "time"
import "database/sql"


type Payment struct {
    DueDate    time.Time      `db:"due_date" json:"due_date"`
    Amount     float64        `db:"amount" json:"amount"`
    PaidAmount sql.NullFloat64 `db:"paid_amount" json:"paid_amount"` 
    PaidAt     *time.Time     `db:"paid_at" json:"paid_at,omitempty"`
	Status     string     `json:"status"`
}