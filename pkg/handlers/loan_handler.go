package handlers

import (
	"encoding/json"
	"net/http"
	"jusan_demo/pkg/middleware"
	"jusan_demo/pkg/models"
	"jusan_demo/pkg/services"
)

type LoanHandlers struct {
	LoanService   *services.LoanService
	ProfileService *services.ProfileService
}

func NewLoanHandlers(loanService *services.LoanService, profileService *services.ProfileService) *LoanHandlers {
	return &LoanHandlers{
		LoanService:   loanService,
		ProfileService: profileService,
	}
}

// @Summary Получить список кредитов пользователя
// @Description Возвращает все кредиты пользователя с графиком платежей
// @Tags loan
// @Produce json
// @Success 200 {array} models.Loan
// @Failure 400 {string} string "Ошибка получения кредитов"
// @Router /api/loans [get]
// @Security BearerAuth
func (h *LoanHandlers) MakeGetLoansHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user from context
		user, ok := r.Context().Value(middleware.UserKey).(*models.AuthUser)
		if !ok {
			http.Error(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		// Get person ID from profile
		profile, err := h.ProfileService.GetProfile(user.ID)
		if err != nil || profile.PersonID == nil {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}

		loans, err := h.LoanService.GetLoans(*profile.PersonID)
		if err != nil {
			http.Error(w, "Error getting loans: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(loans)
	}
}


// @Summary Создать кредит
// @Description Оформляет новый кредит. Пользователи могут создавать только для себя, админы - для любого.
// @Tags loan
// @Accept json
// @Produce json
// @Param request body services.CreateLoanRequest true "Данные кредита"
// @Success 200 {object} services.CreateLoanResponse
// @Failure 400 {string} string "Ошибка создания кредита"
// @Failure 403 {string} string "Недостаточно прав"
// @Router /api/loan/create [post]
// @Security BearerAuth
func (h *LoanHandlers) MakeCreateLoanHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        var req services.CreateLoanRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
            return
        }

        // Get user from context
        user, ok := r.Context().Value(middleware.UserKey).(*models.AuthUser)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // For regular users, ensure they're creating loan for themselves
        if user.Role == "user" {
            profile, err := h.ProfileService.GetProfile(user.ID)
            if err != nil || profile.PersonID == nil {
                http.Error(w, "Profile not found", http.StatusForbidden)
                return
            }
            
            if req.PersonID != *profile.PersonID {
                http.Error(w, "Can only create loan for yourself", http.StatusForbidden)
                return
            }
        }

        // For admins, they can create loans for anyone
        resp, err := h.LoanService.CreateLoan(req)
        if err != nil {
            http.Error(w, "Error creating loan: "+err.Error(), http.StatusBadRequest)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }
}


func MakeUpdateLoanHandler(service *services.LoanService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req services.UpdateLoanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		err := service.UpdateLoan(req)
		if err != nil {
			http.Error(w, "Ошибка обновления кредита: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Кредит успешно обновлён"))
	})
}



type RepayLoanRequest struct {
	LoanID     int     `json:"loan_id"`
	Amount     float64 `json:"amount"`
	RedactorID int     `json:"redactor_id"`
}




//Досрочное погашение 
// @Summary Досрочное погашение кредита
// @Description Выплата части или всего кредита с обновлением графика
// @Tags loan
// @Accept json
// @Produce json
// @Param request body services.RepayLoanRequest true "Данные погашения"
// @Success 200 {string} string "Успешно"
// @Failure 400 {string} string "Ошибка"
// @Router /api/loan/repay [post]
// @Security BearerAuth
func MakeRepayLoanHandler(service *services.LoanService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req services.RepayLoanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		msg, err := service.RepayLoan(req)
		if err != nil {
			http.Error(w, "Ошибка погашения: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	}
}