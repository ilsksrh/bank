package services

import (
	"jusan_demo/pkg/db"
)

type CardService struct{}

func NewCardService() *CardService {
	return &CardService{}
}

type UpsertCardRequest struct {
	CardNumber     string  `json:"card_number"`
	ExpirationDate string  `json:"expiration_date" `
	TypeCode       string  `json:"type_code"`
	StateCode      string  `json:"state_code"`
	PersonID       int     `json:"person_id"`
	AccountID      int     `json:"account_id"`
	IssueDate      *string `json:"issue_date,omitempty"`
	IsDefault      bool    `json:"is_default"`
	IsVirtual      bool    `json:"is_virtual"`
	CardID         *int    `json:"card_id"`
	RedactorID     *int    `json:"redactor_id"`
}

type UpsertCardResponse struct {
	CardID     int
	CardNumber string
}

func (s *CardService) UpsertCard(req UpsertCardRequest) (*UpsertCardResponse, error) {
	var resp UpsertCardResponse
	err := db.DB.QueryRowx(`select * from fn_upsert_card(
		$1, $2, $3, $4, $5, $6,
		$7, $8, $9, $10, $11
	)`,
		req.CardNumber, req.ExpirationDate, req.TypeCode, req.StateCode,
		req.PersonID, req.AccountID,
		req.IssueDate, req.IsDefault, req.IsVirtual,
		req.CardID, req.RedactorID,
	).StructScan(&resp)

	if err != nil {
		return nil, err
	}
	return &resp, nil
}
