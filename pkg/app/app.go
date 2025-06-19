package app

import (
	"jusan_demo/pkg/auth"
	"jusan_demo/pkg/handlers"
	"jusan_demo/pkg/middleware"
	"jusan_demo/pkg/services"
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Инициализируем все сервисы централизованно
	appServices := services.NewAppServices()

	authHandler := auth.NewAuthHandler(appServices.AuthService)

	// Роуты
	mux.Handle("/api/account/upsert",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(handlers.MakeUpsertAccountHandler(appServices.AccountService)),
		),
	)

	mux.Handle("/api/branch/upsert",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(handlers.MakeUpsertBranchHandler(appServices.BranchService)),
		),
	)

	mux.Handle("/api/card/upsert",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(handlers.MakeUpsertCardHandler(appServices.CardService)),
		),
	)

	mux.Handle("/api/person/upsert",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(handlers.MakeUpsertPersonHandler(appServices.PersonService)),
		),
	)

	mux.Handle("/api/transaction/create",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(handlers.MakeCreateTransactionHandler(appServices.TransactionService)),
		),
	)

	mux.Handle("/api/product/upsert",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(handlers.MakeUpsertProductHandler(appServices.ProductService)),
		),
	)

	mux.Handle("/api/loan/create",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(handlers.MakeCreateLoanHandler(appServices.LoanService)),
		),
	)

	mux.Handle("/api/loan/update",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(handlers.MakeUpdateLoanHandler(appServices.LoanService)),
		),
	)

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	mux.Handle("/api/profile",
		middleware.AuthMiddleware(appServices.AuthService, "user")(
			http.HandlerFunc(authHandler.GetProfile)),
	)

	return mux
}
