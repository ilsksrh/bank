package services

import (
	"jusan_demo/pkg/db"
)

type ProductService struct{}

func NewProductService() *ProductService {
	return &ProductService{}
}

type Product struct {
	ID             int     `db:"id" json:"id"`
	Name           string  `db:"name" json:"name"`
	Code           string  `db:"code" json:"code"`
	Description    string  `db:"description" json:"description"`
	InterestRate   float64 `db:"interest_rate" json:"interest_rate"`
	MinAmount      float64 `db:"min_amount" json:"min_amount"`
	MaxAmount      float64 `db:"max_amount" json:"max_amount"`
	TermMinMonths  int     `db:"term_min_months" json:"term_min_months"`
	TermMaxMonths  int     `db:"term_max_months" json:"term_max_months"`
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

func (s *ProductService) GetAllProducts() ([]Product, error) {
	var products []Product
	err := db.DB.Select(&products, `
        SELECT id, name, code, description, interest_rate, min_amount, max_amount, 
               term_min_months, term_max_months 
        FROM product
        WHERE state_id IN (SELECT id FROM states WHERE code = 'PR1') -- только активные
    `)
	return products, err
}