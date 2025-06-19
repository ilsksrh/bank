package services

import (
	"jusan_demo/pkg/db"
)

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

type UpsertAccountRequest struct {
	PersonID       int     `json:"person_id"`
	AccountID      *int    `json:"account_id"`
	TypeCode       *string `json:"type_code"`
	StateCode      *string `json:"state_code"`
	AccountBalance float64 `json:"account_balance"`
	RedactorID     *int    `json:"redactor_id"`
}

type UpsertAccountResponse struct {
	AccountID     int    `db:"account_id"`    
	AccountNumber string `db:"account_number"`
}

func (s *AccountService) UpsertAccount(req UpsertAccountRequest) (*UpsertAccountResponse, error) {
	var resp UpsertAccountResponse
	err := db.DB.QueryRowx(`select * from fn_upsert_account($1, $2, $3, $4, $5, $6)`,
		req.PersonID,
		req.AccountID,
		req.TypeCode,
		req.StateCode,
		req.AccountBalance,
		req.RedactorID,
	).StructScan(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
