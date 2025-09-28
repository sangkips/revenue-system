package user

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/sangkips/revenue-system/internal/domain/user/models"
)

type Handler struct {
	svc *Service
}

func NewHandler(db models.DBTX) *Handler {
	repo := NewRepository(db)

	return &Handler{svc: NewService(repo)}
}

func (h Handler) RegisterUserRoutes(r chi.Router) {
	r.Post("/", h.CreateUser)
	// r.Get("/{id}", h.GetUser)
	// r.Get("/", h.ListUsers)
	// r.Put("/{id}", h.UpdateUser)
	// r.Put("/{id}/password", h.UpdatePassword)
	// r.Delete("/{id}", h.DeleteUser)
}

func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	user, err := h.svc.CreateUser(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
