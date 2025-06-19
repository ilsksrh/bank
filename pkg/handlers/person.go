package handlers

import (
	"encoding/json"
	"jusan_demo/pkg/services"
	"net/http"
)

func MakeUpsertPersonHandler(service *services.PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		userIDCtx := r.Context().Value("userID")
		if userIDCtx == nil {
			http.Error(w, "Не найден userID в контексте", http.StatusUnauthorized)
			return
		}
		userID, ok := userIDCtx.(int)
		if !ok {
			http.Error(w, "Неверный формат userID", http.StatusUnauthorized)
			return
		}

		var req services.UpsertPersonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		if err := service.UpsertPerson(req, userID); err != nil {
			http.Error(w, "Ошибка: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Зарегистрирован"))
	}
}
