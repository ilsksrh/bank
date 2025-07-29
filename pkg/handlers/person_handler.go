package handlers

import (
	"encoding/json"
	"jusan_demo/pkg/services"
	"jusan_demo/pkg/middleware"
	"net/http"
)

func MakeUpsertPersonHandler(service *services.PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Extract userID from context
		userIDCtx := r.Context().Value(middleware.UserIDKey)

		if userIDCtx == nil {
			http.Error(w, "Не найден userID в контексте", http.StatusUnauthorized)
			return
		}
		userID, ok := userIDCtx.(int)
		if !ok {
			http.Error(w, "Неверный формат userID", http.StatusUnauthorized)
			return
		}

		// Extract role from context
		roleCtx := r.Context().Value(middleware.RoleKey)

		if roleCtx == nil {
			http.Error(w, "Не найдена роль пользователя", http.StatusUnauthorized)
			return
		}
		role, ok := roleCtx.(string)
		if !ok {
			http.Error(w, "Неверный формат роли", http.StatusUnauthorized)
			return
		}
		isAdmin := role == "admin"

		var req services.UpsertPersonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		if err := service.UpsertPerson(req, userID, isAdmin); err != nil {
			http.Error(w, "Ошибка: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":   true,
			"message":   "Зарегистрирован",
		})
	}
}