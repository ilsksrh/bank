package handlers

import (
	"encoding/json"
	"jusan_demo/pkg/services"
	"net/http"
)

// @Summary Обновить или создать продукт
// @Description Добавляет или обновляет продукт банка (кредит, услуга и т.п.)
// @Tags product
// @Accept json
// @Produce json
// @Param request body services.UpsertProductRequest true "Данные продукта"
// @Success 200 {object} services.UpsertProductResponse
// @Failure 400 {string} string "Ошибка при сохранении"
// @Router /api/product/upsert [post]
// @Security BearerAuth
func MakeUpsertProductHandler(service *services.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var req services.UpsertProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := service.UpsertProduct(req)
		if err != nil {
			http.Error(w, "Ошибка: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// @Summary Получить все кредитные продукты
// @Description Возвращает список кредитных продуктов
// @Tags product
// @Produce json
// @Success 200 {array} services.Product
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/product/list [get]
func MakeListProductHandler(service *services.ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := service.GetAllProducts()
		if err != nil {
			http.Error(w, "Ошибка получения продуктов: "+err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(products)
	}
}
