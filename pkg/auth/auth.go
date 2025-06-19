package auth


import (
	"encoding/json"
	"net/http"
	"log"
)

type AuthHandler struct {
	Service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}


type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	PersonID *int    `json:"person_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ProfileResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	PersonID int    `json:"person_id"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Неверный JSON в /register:", err)
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	pid := 0
	if req.PersonID != nil {
		pid = *req.PersonID
	}

	err := h.Service.Register(req.Email, req.Password, req.Role, &pid)
	if err != nil {
		if err.Error() == "email уже используется" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			log.Println("Ошибка регистрации:", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Регистрация успешна: %s", req.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Регистрация прошла успешно",
	})
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Неверный JSON в /login:", err)
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	access, refresh, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		log.Println("Ошибка логина:", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Printf("Успешный вход: %s", req.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
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

	profile, err := h.Service.GetProfile(userID)
	if err != nil {
		log.Println("Ошибка получения профиля:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if profile.PersonID == nil {
		log.Printf("Профиль не завершён для пользователя %d", userID)
		http.Error(w, "Пожалуйста, завершите регистрацию, заполнив профиль", http.StatusUnprocessableEntity)
		return
	}

	log.Printf("Профиль успешно получен: %s", profile.Email)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

