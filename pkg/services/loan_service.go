package services

import (
	"jusan_demo/pkg/db"
)

type LoanService struct{}

func NewLoanService() *LoanService {
	return &LoanService{}
}

type CreateLoanRequest struct {
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
	TermMonths   int     `json:"term_months"`
	StartDate    string  `json:"start_date"`
	PersonID     int     `json:"person_id"`
	ProductID    int     `json:"product_id"`
	StateCode    string  `json:"state_code"`
	RedactorID   int     `json:"redactor_id"`
}
type CreateLoanResponse struct {
	LoanID int `db:"loan_id"`
}

func (s *LoanService) CreateLoan(req CreateLoanRequest) (CreateLoanResponse, error) {
	var resp CreateLoanResponse
	err := db.DB.QueryRowx(`select * from fn_create_loan(
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)`,
		req.Name, req.Amount, req.InterestRate, req.TermMonths,
		req.StartDate, req.PersonID, req.ProductID, req.StateCode, req.RedactorID,
	).StructScan(&resp)
	return resp, err
}


type UpdateLoanRequest struct {
	LoanID       int     `json:"loan_id"`
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
	TermMonths   int     `json:"term_months"`
	StartDate    string  `json:"start_date"`
	PersonID     int     `json:"person_id"`
	ProductID    int     `json:"product_id"`
	StateCode    string  `json:"state_code"`
	RedactorID   int     `json:"redactor_id"`
}

func (s *LoanService) UpdateLoan(req UpdateLoanRequest) error {
	_, err := db.DB.Exec(`call update_loan(
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	)`,
		req.LoanID, req.Name, req.Amount, req.InterestRate,
		req.TermMonths, req.StartDate, req.PersonID,
		req.ProductID, req.StateCode, req.RedactorID,
	)
	return err
}
