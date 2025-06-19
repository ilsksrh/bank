	package handlers

	import (
		"encoding/json"
		"jusan_demo/pkg/services"
		"net/http"
	)

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
