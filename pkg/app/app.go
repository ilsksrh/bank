package app

import (
	"encoding/json"
	"net/http"

	"jusan_demo/pkg/auth"
	"jusan_demo/pkg/handlers"
	"jusan_demo/pkg/middleware"
	"jusan_demo/pkg/models"
	"jusan_demo/pkg/services"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "jusan_demo/docs"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize services
	appServices := services.NewAppServices()
	authService := auth.NewAuthService()

	// Auth handler
	authHandler := auth.NewAuthHandler(authService, appServices.ProfileService)

	// Documentation
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("OK")
	})

	// Public routes
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// Admin and user routes
	setupAdminRoutes(mux, authService, appServices)
	setupUserRoutes(mux, authService, appServices)

	// Fallback
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return mux
}

func setupAdminRoutes(mux *http.ServeMux, authService *auth.AuthService, appServices *services.AppServices) {
	routes := []struct {
		path    string
		handler http.Handler
	}{
		{"/api/account/upsert", handlers.MakeUpsertAccountHandler(appServices.AccountService)},
		{"/api/branch/upsert", handlers.MakeUpsertBranchHandler(appServices.BranchService)},
		{"/api/card/upsert", handlers.MakeUpsertCardHandler(appServices.CardService)},
		{"/api/product/upsert", handlers.MakeUpsertProductHandler(appServices.ProductService)},
		{"/api/loan/update", handlers.MakeUpdateLoanHandler(appServices.LoanService)},
	}

	for _, route := range routes {
		mux.Handle(route.path,
			middleware.AuthMiddleware(authService, "admin")(route.handler),
		)
	}
}

func setupUserRoutes(mux *http.ServeMux, authService *auth.AuthService, appServices *services.AppServices) {
	// LoanHandlers
	loanHandlers := handlers.NewLoanHandlers(appServices.LoanService, appServices.ProfileService)

	// Обёртка с проверкой авторизации
	authWrapper := func(h http.HandlerFunc) http.Handler {
		return middleware.AuthMiddleware(authService, "user", "admin")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, ok := r.Context().Value(middleware.UserKey).(*models.AuthUser); !ok {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				h(w, r)
			}),
		)
	}

	// Routes
	mux.Handle("/api/profile", authWrapper(
		handlers.MakeGetProfileHandler(appServices.ProfileService),
	))

	mux.Handle("/api/product/list", authWrapper(
		handlers.MakeListProductHandler(appServices.ProductService),
	))

	mux.Handle("/api/loan/create",
		middleware.AuthMiddleware(authService, "user", "admin")(
			loanHandlers.MakeCreateLoanHandler(),
		),
	)

	mux.Handle("/api/loans",
		middleware.AuthMiddleware(authService, "user", "admin")(
			loanHandlers.MakeGetLoansHandler(),
		),
	)
	mux.Handle("/api/person/upsert", authWrapper(
		handlers.MakeUpsertPersonHandler(appServices.PersonService),
	))

	mux.Handle("/api/transaction/create", authWrapper(
		handlers.MakeCreateTransactionHandler(appServices.TransactionService),
	))

	mux.Handle("/api/loan/repay", authWrapper(
	handlers.MakeRepayLoanHandler(appServices.LoanService),
))

}
