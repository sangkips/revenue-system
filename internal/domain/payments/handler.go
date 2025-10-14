package payments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
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

	// Payment Allocations sub-routes
	r.Route("/{id}/allocations", func(r chi.Router) {
		r.Post("/", h.CreatePaymentAllocation)
		r.Get("/", h.ListPaymentAllocations)
		r.Delete("/{allocation_id}", h.DeletePaymentAllocation)
	})

	// Receipts sub-routes
	r.Route("/{id}/receipts", func(r chi.Router) {
		r.Post("/", h.CreateReceipt)
		r.Get("/{receipt_id}", h.GetReceipt)
		r.Get("/", h.ListReceiptsByPayment)
		r.Patch("/{receipt_id}", h.UpdateReceipt)
		r.Delete("/{receipt_id}", h.DeleteReceipt)
	})
}

func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest 

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

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
	if id == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}
	
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
	if id == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}
	
	var req UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	ctx := r.Context()

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
	if id == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}
	
	ctx := r.Context()
	if err := h.svc.DeletePayment(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


// Payment Allocations Handlers
func (h *Handler) CreatePaymentAllocation(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	var req CreatePaymentAllocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.PaymentID = paymentID
	ctx := r.Context()
	allocation, err := h.svc.CreatePaymentAllocation(ctx, req)
	if err != nil {
		log.Error().Err(err).Str("payment_id", paymentID).Msg("Failed to create allocation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(allocation)
}

func (h *Handler) ListPaymentAllocations(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	log.Info().Str("payment_id", paymentID).Msg("Listing payment allocations")
	
	if paymentID == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}
	
	ctx := r.Context()
	allocations, err := h.svc.ListPaymentAllocations(ctx, paymentID)
	if err != nil {
		log.Error().Err(err).Str("payment_id", paymentID).Msg("Failed to list payment allocations")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.Info().Int("count", len(allocations)).Str("payment_id", paymentID).Msg("Found payment allocations")
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allocations)
}

func (h *Handler) DeletePaymentAllocation(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	allocationID := chi.URLParam(r, "allocation_id")
	ctx := r.Context()
	if err := h.svc.DeletePaymentAllocation(ctx, allocationID, paymentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Receipts Handlers
func (h *Handler) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	var req CreateReceiptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.PaymentID = paymentID
	ctx := r.Context()
	if err := h.svc.CreateReceipt(ctx, req); err != nil {
		log.Error().Err(err).Str("payment_id", paymentID).Msg("Failed to create receipt")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetReceipt(w http.ResponseWriter, r *http.Request) {
	receiptID := chi.URLParam(r, "receipt_id")
	ctx := r.Context()
	receipt, err := h.svc.GetReceipt(ctx, receiptID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(receipt)
}

func (h *Handler) ListReceiptsByPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	ctx := r.Context()
	receipts, err := h.svc.ListReceiptsByPayment(ctx, paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(receipts)
}

func (h *Handler) UpdateReceipt(w http.ResponseWriter, r *http.Request) {
	receiptID := chi.URLParam(r, "receipt_id")
	var req UpdateReceiptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if err := h.svc.UpdateReceipt(ctx, receiptID, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteReceipt(w http.ResponseWriter, r *http.Request) {
	receiptID := chi.URLParam(r, "receipt_id")
	ctx := r.Context()
	if err := h.svc.DeleteReceipt(ctx, receiptID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
