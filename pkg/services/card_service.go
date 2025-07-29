package services

import (
	//"errors"
	"fmt"
	"jusan_demo/pkg/db"
	"time"
)

type CardService struct{}

func NewCardService() *CardService {
	return &CardService{}
}

type UpsertCardRequest struct {
	CardNumber     *string `json:"card_number,omitempty"`
	ExpirationDate string  `json:"expiration_date"`
	TypeCode       string  `json:"type_code"`
	StateCode      string  `json:"state_code"`
	PersonID       int     `json:"person_id"`
	AccountID      int     `json:"account_id"`
	IssueDate      *string `json:"issue_date,omitempty"`
	IsDefault      bool    `json:"is_default"`
	IsVirtual      bool    `json:"is_virtual"`
	CardID         *int    `json:"card_id,omitempty"`
	RedactorID     *int    `json:"redactor_id,omitempty"`
}

type UpsertCardResponse struct {
	CardID     int    `db:"card_id"`
	CardNumber string `db:"card_number"`
}

func (s *CardService) UpsertCard(req UpsertCardRequest) (*UpsertCardResponse, error) {
    var resp UpsertCardResponse

    // Parse expiration date
    expirationDate, err := time.Parse("2006-01-02", req.ExpirationDate)
    if err != nil {
        return nil, fmt.Errorf("invalid expiration date format, expected YYYY-MM-DD: %w", err)
    }

    // Handle issue date
    var issueDate time.Time
    if req.IssueDate != nil {
        issueDate, err = time.Parse("2006-01-02", *req.IssueDate)
        if err != nil {
            return nil, fmt.Errorf("invalid issue date format, expected YYYY-MM-DD: %w", err)
        }
    } else {
        issueDate = time.Now()
    }

    // Convert optional fields to appropriate SQL types
    var cardNumber interface{} = req.CardNumber
    if req.CardNumber == nil {
        cardNumber = nil
    }

    var cardID interface{} = req.CardID
    if req.CardID == nil {
        cardID = nil
    }

    var redactorID interface{} = req.RedactorID
    if req.RedactorID == nil {
        redactorID = nil
    }

    err = db.DB.QueryRowx(`
        SELECT * FROM fn_upsert_card(
            $1::date,    -- expiration_date
            $2,          -- type_code
            $3,          -- state_code
            $4,          -- person_id
            $5,          -- account_id
            $6,          -- card_number
            $7::date,    -- issue_date
            $8,          -- is_default
            $9,          -- is_virtual
            $10,         -- card_id
            $11          -- redactor_id
        )`,
        expirationDate, // $1
        req.TypeCode,   // $2
        req.StateCode,  // $3
        req.PersonID,   // $4
        req.AccountID,  // $5
        cardNumber,     // $6
        issueDate,      // $7
        req.IsDefault,  // $8
        req.IsVirtual,  // $9
        cardID,         // $10
        redactorID,     // $11
    ).StructScan(&resp)

    if err != nil {
        return nil, fmt.Errorf("error upserting card: %w", err)
    }
    return &resp, nil
}