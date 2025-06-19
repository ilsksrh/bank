package services

import (
	"jusan_demo/pkg/db"
)

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	TypeCode    string  `json:"type_code"`
	StateCode   string  `json:"state_code"`
	PersonID    int     `json:"person_id"`
	RedactorID  int     `json:"redactor_id"`
	AccountID   *int    `json:"account_id"`
	CardID      *int    `json:"card_id"`
	LoanID      *int    `json:"loan_id"`
}

type CreateTransactionResponse struct {
	TransactionID int `db:"transaction_id"`
}

func (s *TransactionService) CreateTransaction(req CreateTransactionRequest) (CreateTransactionResponse, error) {
	var resp CreateTransactionResponse
	err := db.DB.QueryRowx(`select * from fn_create_transaction(
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)`,
		req.Amount, req.Description, req.TypeCode, req.StateCode,
		req.PersonID, req.RedactorID, req.AccountID, req.CardID, req.LoanID,
	).StructScan(&resp)
	return resp, err
}
