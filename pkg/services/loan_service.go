package services

import (
	"fmt"
	"math"
	"time"
	"jusan_demo/pkg/db"
	"jusan_demo/pkg/models"
)

type LoanService struct{}

func NewLoanService() *LoanService {
	return &LoanService{}
}

type CreateLoanRequest struct {
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	TermMonths int     `json:"term_months"`
	PersonID   int     `json:"person_id"`
	ProductID  int     `json:"product_id"`
	StateCode  string  `json:"state_code"`
	RedactorID *int    `json:"redactor_id,omitempty"`
}


type CreateLoanResponse struct {
	LoanID          int              `json:"loan_id"`
	MonthlyPayment  float64          `json:"monthly_payment"`
	TotalPayment    float64          `json:"total_payment"`
	Overpayment     float64          `json:"overpayment"`
	PaymentSchedule []models.Payment `json:"payment_schedule"`
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

func (s *LoanService) CreateLoan(req CreateLoanRequest) (*CreateLoanResponse, error) {
	var product models.Product
	if err := db.DB.Get(&product, `
        SELECT id, name, interest_rate, min_amount, max_amount, 
               term_min_months, term_max_months
        FROM product 
        WHERE id = $1`, req.ProductID); err != nil {
		return nil, fmt.Errorf("invalid product: %w", err)
	}

	if req.Amount < product.MinAmount || req.Amount > product.MaxAmount {
		return nil, fmt.Errorf("amount must be between %.2f and %.2f",
			product.MinAmount, product.MaxAmount)
	}

	if req.TermMonths < product.TermMinMonths || req.TermMonths > product.TermMaxMonths {
		return nil, fmt.Errorf("term must be between %d and %d months",
			product.TermMinMonths, product.TermMaxMonths)
	}

	rate := product.InterestRate
	var redactorIDPtr *int
	if req.RedactorID != nil {
		redactorIDPtr = req.RedactorID
	}

	var loanID int
	startDate := time.Now()

	if err := db.DB.Get(&loanID, `
		SELECT loan_id FROM fn_create_loan(
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`,
		req.Name, req.Amount, rate, req.TermMonths,
		startDate, req.PersonID, req.ProductID,
		req.StateCode, redactorIDPtr,
	); err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	r := rate / 12 / 100
	var monthly float64
	if r == 0 {
		monthly = math.Round(req.Amount/float64(req.TermMonths)*100) / 100
	} else {
		pow := math.Pow(1+r, float64(req.TermMonths))
		monthly = math.Round((req.Amount*r*pow/(pow-1))*100) / 100
	}

	var schedule []models.Payment
	if err := db.DB.Select(&schedule, `
		SELECT due_date, amount, paid_amount, paid_at
		FROM vw_payment_schedule 
		WHERE loan_id = $1 
		ORDER BY due_date`, loanID); err != nil {
		return nil, fmt.Errorf("failed to get payment schedule: %w", err)
	}

	total := monthly * float64(req.TermMonths)
	return &CreateLoanResponse{
		LoanID:          loanID,
		MonthlyPayment:  monthly,
		TotalPayment:    total,
		Overpayment:     total - req.Amount,
		PaymentSchedule: schedule,
	}, nil
}



func (s *LoanService) RepayLoanEarly(loanID int, amount float64, redactorID int) error {
	var result string
	err := db.DB.Get(&result, `
		SELECT fn_repay_loan_early($1, $2, $3)
	`, loanID, amount, redactorID)
	if err != nil {
		return fmt.Errorf("ошибка досрочного погашения: %w", err)
	}
	return nil
}



func (s *LoanService) UpdateLoan(req UpdateLoanRequest) error {
	_, err := db.DB.Exec(`
		SELECT fn_update_loan(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`,
		req.LoanID, req.Name, req.Amount, req.InterestRate,
		req.TermMonths, req.StartDate, req.PersonID,
		req.ProductID, req.StateCode, req.RedactorID,
	)
	return err
}

func (s *LoanService) GetLoans(personID int) ([]models.Loan, error) {
	var loans []models.Loan
	err := db.DB.Select(&loans, `
        SELECT 
            l.id, 
            l.name, 
            l.amount, 
            l.remaining_amount,
            l.interest_rate, 
            l.term_months, 
            l.start_date, 
            l.end_date,
            s.code as state_code,
            p.name as product_name
        FROM loan l
        JOIN states s ON l.state_id = s.id
        JOIN product p ON l.product_id = p.id
        WHERE l.person_id = $1
        ORDER BY l.start_date DESC`, personID)

	if err != nil {
		return nil, fmt.Errorf("failed to get loans: %w", err)
	}

	for i := range loans {
		var payments []models.Payment
		err = db.DB.Select(&payments, `
            SELECT 
                due_date, 
                amount,
                COALESCE(paid_amount, 0) as paid_amount,
                paid_at
            FROM payment_schedule
            WHERE loan_id = $1
            ORDER BY due_date`, loans[i].ID)

		if err != nil {
			err = db.DB.Select(&payments, `
                SELECT 
                    due_date, 
                    amount,
                    0 as paid_amount,
                    NULL as paid_at
                FROM payment_schedule
                WHERE loan_id = $1
                ORDER BY due_date`, loans[i].ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get payment schedule: %w", err)
			}
		}

		loans[i].PaymentSchedule = payments
	}

	return loans, nil
}