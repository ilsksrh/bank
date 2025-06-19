package services

import "jusan_demo/pkg/auth"

type AppServices struct {
	AuthService        *auth.AuthService
	PersonService      *PersonService
	AccountService     *AccountService
	LoanService        *LoanService
	TransactionService *TransactionService
	CardService        *CardService
	BranchService      *BranchService
	ProductService     *ProductService
}

func NewAppServices() *AppServices {
	return &AppServices{
		AuthService:        auth.NewAuthService(),
		PersonService:      NewPersonService(),
		AccountService:     NewAccountService(),
		LoanService:        NewLoanService(),
		TransactionService: NewTransactionService(),
		CardService:        NewCardService(),
		BranchService:      NewBranchService(),
		ProductService:     NewProductService(),
	}
}
