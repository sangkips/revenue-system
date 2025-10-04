package counties

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/sangkips/revenue-system/internal/domain/counties/models"
)

type Handler struct {
	svc *Service
}

func NewHandler(db models.DBTX) *Handler {
	repo := NewRepository(db)
	return &Handler{svc: NewService(repo)}
}

func (h *Handler) RegisterCountyRoutes(r chi.Router) {
	r.Post("/", h.CreateCounty)
	r.Get("/{id}", h.GetCounty)
	r.Get("/", h.ListCounties)
	r.Patch("/{id}", h.UpdateCounty)
	r.Delete("/{id}", h.DeleteCounty)
}

func (h *Handler) CreateCounty(w http.ResponseWriter, r *http.Request) {
	var req CreateCountyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	county, err := h.svc.CreateCounty(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create county")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(county)
}

func (h *Handler) GetCounty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	county, err := h.svc.GetCounty(ctx, int32(id))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(county)
}

func (h *Handler) ListCounties(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10,32)

	if limit == 0 {
		limit = 10
	}

	ctx := r.Context()
	county, err := h.svc.ListCounties(ctx, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(county)
}

func (h *Handler) UpdateCounty(w http.ResponseWriter, r *http.Request) {
	countyIDStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(countyIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateCountyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	county, err := h.svc.UpdateCounty(ctx, int32(id), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(county)
}

func (h *Handler) DeleteCounty(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.svc.DeleteCounty(ctx, int32(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
