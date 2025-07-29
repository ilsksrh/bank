package models

import (
	"time"
)

type Person struct {
	ID          *int      `db:"id" json:"id,omitempty"` 
	FirstName   string    `db:"first_name" json:"first_name" binding:"required"`
	LastName    string    `db:"last_name" json:"last_name" binding:"required"`
	DateOfBirth string    `db:"date_of_birth" json:"date_of_birth" binding:"required"` 
	UIN         string    `db:"uin" json:"uin" binding:"required"`
	PhoneNumber string    `db:"phone_number" json:"phone_number" binding:"required"`
	StateID     *int      `db:"state_id" json:"state_id,omitempty"`
	BranchID    *int      `db:"branch_id" json:"branch_id,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at,omitempty"`
	Accounts    []Account `json:"accounts,omitempty"`
}