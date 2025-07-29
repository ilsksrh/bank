package services



type AppServices struct {
	PersonService      *PersonService
	AccountService     *AccountService
	LoanService       *LoanService
	ProfileService    *ProfileService
	TransactionService *TransactionService
	CardService        *CardService
	BranchService      *BranchService
	ProductService     *ProductService
}

func NewAppServices() *AppServices {
	return &AppServices{
		PersonService:      NewPersonService(),
		AccountService:     NewAccountService(),
		LoanService:        NewLoanService(),
		TransactionService: NewTransactionService(),
		CardService:        NewCardService(),
		BranchService:      NewBranchService(),
		ProductService:     NewProductService(),
	}
}