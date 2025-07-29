package handlers

import (
	"encoding/json"
	"fmt"
	"jusan_demo/pkg/services"
	"net/http"
)
// @Summary Создать транзакцию
// @Description Создаёт финансовую транзакцию по счёту/карте/кредиту
// @Tags transaction
// @Accept json
// @Produce json
// @Param request body services.CreateTransactionRequest true "Данные транзакции"
// @Success 200 {object} services.CreateTransactionResponse
// @Failure 400 {string} string "Ошибка транзакции"
// @Router /api/transaction/create [post]
// @Security BearerAuth
func MakeCreateTransactionHandler(service *services.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req services.CreateTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.CreateTransaction(req)
		if err != nil {
			http.Error(w, "Ошибка: "+err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("Тело запроса: %+v\n", req)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
