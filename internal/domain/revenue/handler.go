package revenue

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/sangkips/revenue-system/internal/domain/revenue/models"
)

type Handler struct {
	svc *Service
}

func NewHandler(db models.DBTX) *Handler {
	repo := NewRepository(db)
	return &Handler{svc: NewService(repo)}
}

func (h *Handler) RegisterRevenueRoutes(r chi.Router) {
	r.Post("/", h.CreateRevenue)
	r.Get("/{id}", h.GetRevenue)
	r.Get("/", h.ListRevenues)
	r.Get("/taxpayer/{taxpayer_id}", h.ListRevenuesByTaxpayerID)
	r.Patch("/{id}", h.UpdateRevenue)
	r.Delete("/{id}", h.DeleteRevenue)
}

func (h *Handler) CreateRevenue(w http.ResponseWriter, r *http.Request) {
	var req CreateRevenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	revenue, err := h.svc.CreateRevenue(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create revenue")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(revenue)
}

func (h *Handler) GetRevenue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	revenue, err := h.svc.GetRevenue(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(revenue)
}

func (h *Handler) ListRevenues(w http.ResponseWriter, r *http.Request) {
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
	revenue, err := h.svc.ListRevenues(ctx, int32(countyID), int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(revenue)
}

func (h *Handler) UpdateRevenue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateRevenueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	revenue, err := h.svc.UpdateRevenue(ctx, id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(revenue)
}

func (h *Handler) ListRevenuesByTaxpayerID(w http.ResponseWriter, r *http.Request) {
	taxpayerID := chi.URLParam(r, "taxpayer_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	if limit == 0 {
		limit = 10
	}

	ctx := r.Context()
	revenue, err := h.svc.ListRevenuesByTaxpayerID(ctx, taxpayerID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(revenue)
}

func (h *Handler) DeleteRevenue(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	if err := h.svc.DeleteRevenue(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
