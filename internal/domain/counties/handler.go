package counties

import (
	"encoding/json"
	"net/http"

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
}

func (h *Handler) CreateCounty(w http.ResponseWriter, r *http.Request) {
	var req CreateCountyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if err := h.svc.CreateCounty(ctx, req); err != nil {
		log.Error().Err(err).Msg("Failed to create county")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
