package handlers

import (
	"encoding/json"
	"jusan_demo/pkg/services"
	"net/http"
)
// @Summary Обновить или создать карту
// @Description Создаёт или обновляет карту
// @Tags card
// @Accept json
// @Produce json
// @Param request body services.UpsertCardRequest true "Данные карты"
// @Success 200 {object} services.UpsertCardResponse
// @Failure 400 {string} string "Ошибка при создании карты"
// @Router /api/card/upsert [post]
// @Security BearerAuth
func MakeUpsertCardHandler(service *services.CardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req services.UpsertCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.UpsertCard(req)
		if err != nil {
			http.Error(w, "Ошибка: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
