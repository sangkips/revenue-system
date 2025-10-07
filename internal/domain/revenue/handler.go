package revenue

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	r.Patch("/{id}", h.UpdateRevenue)
	r.Delete("/{id}", h.DeleteRevenue)
}

func (h *Handler) CreateRevenue(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetRevenue(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) ListRevenues(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) UpdateRevenue(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteRevenue(w http.ResponseWriter, r *http.Request) {}