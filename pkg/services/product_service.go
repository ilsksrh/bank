package services

import (
	"jusan_demo/pkg/db"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

type UpsertProductRequest struct {
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	InterestRate   float64  `json:"interest_rate"`
	MinAmount      float64  `json:"min_amount"`
	MaxAmount      float64  `json:"max_amount"`
	TermMinMonths  int      `json:"term_min_months"`
	TermMaxMonths  int      `json:"term_max_months"`
	TypeCode       string   `json:"type_code"`
	StateCode      string   `json:"state_code"`
	RedactorID     *int     `json:"redactor_id"`
}

type UpsertProductResponse struct {
	ProductID int    `db:"product_id"`
	Code      string `db:"code"`
}

func (s *ProductService) UpsertProduct(req UpsertProductRequest) (UpsertProductResponse, error) {
	var resp UpsertProductResponse
	err := db.DB.QueryRowx(`select * from fn_upsert_product(
		$1, $2, $3, $4, $5, $6,
		$7, $8, $9, $10, $11
	)`,
		req.Name, req.Code, req.Description,
		req.InterestRate, req.MinAmount, req.MaxAmount,
		req.TermMinMonths, req.TermMaxMonths,
		req.TypeCode, req.StateCode, req.RedactorID,
	).StructScan(&resp)
	return resp, err
}
