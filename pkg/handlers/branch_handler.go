	package handlers

	import (
		"encoding/json"
		"jusan_demo/pkg/services"
		"net/http"
	)
// @Summary Обновить или создать филиал
// @Description Добавляет или обновляет данные филиала
// @Tags branch
// @Accept json
// @Produce json
// @Param request body services.UpsertBranchRequest true "Данные филиала"
// @Success 200 {object} services.UpsertBranchResponse
// @Failure 400 {string} string "Ошибка валидации или запроса"
// @Router /api/branch/upsert [post]
// @Security BearerAuth
	func MakeUpsertBranchHandler(service *services.BranchService) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
				return
			}

			var req services.UpsertBranchRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
				return
			}

			resp, err := service.UpsertBranch(req)
			if err != nil {
				http.Error(w, "Ошибка: "+err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}
	}
