package models

type Card struct {
	ID          int    `db:"id" json:"id"`
	CardNumber  string `db:"card_number" json:"card_number"`
	IsVirtual   bool   `db:"is_virtual" json:"is_virtual"`
	IsDefault   bool   `db:"is_default" json:"is_default"`
	IssueDate   string `db:"issue_date" json:"issue_date"`
	ExpiryDate  string `db:"expiration_date" json:"expiration_date"`
}
