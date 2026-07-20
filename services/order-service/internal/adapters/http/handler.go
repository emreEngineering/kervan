package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/emreEngineering/kervan/services/order-service/internal/application"
)

type OrderHandler struct {
	app *application.OrderService
}

func NewOrderHandler(app *application.OrderService) *OrderHandler {
	return &OrderHandler{app: app}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := extractToken(r)
	if token == "" {
		http.Error(w, "token gerekli", http.StatusUnauthorized)
		return
	}

	order, err := h.app.CreateOrder(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)

}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "geçersiz id", http.StatusBadRequest)
		return
	}

	order, err := h.app.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
func (h *OrderHandler) ReturnOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/orders/")
	idStr = strings.TrimSuffix(idStr, "/return")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "geçersiz id", http.StatusBadRequest)
		return
	}

	order, err := h.app.ReturnOrder(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}
