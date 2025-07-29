package handlers

import (
	"encoding/json"
	"jusan_demo/pkg/services"
	"log"
	"net/http"
)
// @Summary Обновить или создать счёт
// @Description Создаёт новый счёт или обновляет существующий
// @Tags account
// @Accept json
// @Produce json
// @Param request body services.UpsertAccountRequest true "Данные счёта"
// @Success 200 {object} services.UpsertAccountResponse
// @Failure 400 {string} string "Ошибка валидации или запроса"
// @Router /api/account/upsert [post]
// @Security BearerAuth
func MakeUpsertAccountHandler(service *services.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req services.UpsertAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("DEBUG: %+v", req) 

		resp, err := service.UpsertAccount(req)
		if err != nil {
			http.Error(w, "Ошибка: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
