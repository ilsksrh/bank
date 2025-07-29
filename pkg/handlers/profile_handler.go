package handlers

import (
	"encoding/json"
	"net/http"
	"jusan_demo/pkg/models"
	"jusan_demo/pkg/services"
	"jusan_demo/pkg/middleware"
)

// @Summary Get user profile
// @Description Retrieves the authenticated user's profile information including personal details and accounts
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.AuthUser "Successfully retrieved profile"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Profile not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/profile [get]
func MakeGetProfileHandler(service *services.ProfileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.UserKey).(*models.AuthUser)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		profile, err := service.GetProfile(user.ID)
		if err != nil {
			http.Error(w, "Error getting profile: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}