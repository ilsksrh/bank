package services

import (
	"jusan_demo/pkg/db"
)

type BranchService struct{}

func NewBranchService() *BranchService {
	return &BranchService{}
}

type UpsertBranchRequest struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	StateCode    string  `json:"state_code"`
	Location     *string `json:"location"`
	City         *string `json:"city"`
	PhoneNumber  *string `json:"phone_number"`
	RedactorID   *int    `json:"redactor_id"`
}

type UpsertBranchResponse struct {
	BranchID   int    `db:"branch_id"`
	BranchCode string `db:"branch_code"`
}


func (s *BranchService) UpsertBranch(req UpsertBranchRequest) (*UpsertBranchResponse, error) {
	var resp UpsertBranchResponse
	err := db.DB.QueryRowx(`select * from fn_upsert_branch(
		$1, $2, $3, $4, $5, $6, $7
	)`,
		req.Code, req.Name, req.StateCode,
		req.Location, req.City, req.PhoneNumber, req.RedactorID,
	).StructScan(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
