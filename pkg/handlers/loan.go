package handlers

import (
	"encoding/json"
	"jusan_demo/pkg/services"
	"net/http"
)

func MakeCreateLoanHandler(service *services.LoanService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req services.CreateLoanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.CreateLoan(req)
		if err != nil {
			http.Error(w, "Ошибка создания кредита: "+err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func MakeUpdateLoanHandler(service *services.LoanService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req services.UpdateLoanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := service.UpdateLoan(req); err != nil {
			http.Error(w, "Ошибка обновления кредита: "+err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}
}
