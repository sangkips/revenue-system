package payments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sangkips/revenue-system/internal/domain/payments/models"
	"github.com/sangkips/revenue-system/internal/middleware/auth"
)

type Handler struct {
	svc *Service
}

func NewHandler(db models.DBTX) *Handler {
	repo := NewRepository(db)
	return &Handler{svc: NewService(repo)}
}

func (h *Handler) RegisterPaymentsRoutes(r chi.Router) {
	r.Post("/", h.CreatePayment)
	r.Get("/{id}", h.GetPayment)
	r.Get("/", h.ListPayments)
	r.Get("/revenue/{revenue_id}", h.ListPaymentsByRevenueID)
	r.Patch("/{id}", h.UpdatePayment)
	r.Delete("/{id}", h.DeletePayment)
}

func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest 

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	payment, err := h.svc.CreatePayment(ctx, req, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	payment, err := h.svc.GetPayment(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func (h *Handler) ListPayments(w http.ResponseWriter, r *http.Request) {
	countyIDStr := r.URL.Query().Get("county_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	countyID, _ := strconv.ParseInt(countyIDStr, 10, 32)
	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	if limit == 0 {
		limit = 10
	}

	ctx := r.Context()
	payments, err := h.svc.ListPayments(ctx, int32(countyID), int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}

func (h *Handler) ListPaymentsByRevenueID(w http.ResponseWriter, r *http.Request) {
	revenueID := chi.URLParam(r, "revenue_id")
	ctx := r.Context()
	payments, err := h.svc.ListPaymentsByRevenueID(ctx, revenueID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payments)
}

func (h *Handler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	ctx := r.Context()
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	payment, err := h.svc.UpdatePayment(ctx, id, req, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}


func (h *Handler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	if err := h.svc.DeletePayment(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

