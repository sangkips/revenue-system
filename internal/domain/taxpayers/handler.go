package taxpayers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/sangkips/revenue-system/internal/domain/taxpayers/models"
)


type Handler struct {
	svc *Service
}

func NewHandler(db models.DBTX) *Handler {
	repo := NewRepository(db)
	return  &Handler{svc: NewService(repo)}
}

func (h *Handler) RegisterTaxpayerRoutes(r chi.Router) {
	r.Post("/", h.CreateTaxpayer)
	r.Get("/{id}", h.GetTaxpayer)
	r.Get("/", h.ListTaxpayers)
	r.Patch("/{id}", h.UpdateTaxpayer)
	r.Delete("/{id}", h.DeleteTaxpayer)
}

func (h *Handler) CreateTaxpayer(w http.ResponseWriter, r *http.Request) {
	var req CreateTaxpayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	taxpayer, err := h.svc.CreateTaxpayer(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create taxpayer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(taxpayer)
}

func (h *Handler) GetTaxpayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	taxpayer, err := h.svc.GetTaxpayer(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taxpayer)
}

func (h *Handler) ListTaxpayers(w http.ResponseWriter, r *http.Request) {
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

	taxpayer, err := h.svc.ListTaxpayers(ctx, int32(countyID), int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taxpayer)
}

func (h *Handler) UpdateTaxpayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateTaxpayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	taxpayer, err := h.svc.UpdateTaxpayer(ctx, id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taxpayer)
}

func (h *Handler) DeleteTaxpayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	if err := h.svc.DeleteTaxpayer(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}