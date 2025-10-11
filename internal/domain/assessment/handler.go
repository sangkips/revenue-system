package assessment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/sangkips/revenue-system/internal/domain/assessment/models"
	"github.com/sangkips/revenue-system/internal/middleware/auth"
)

type Handler struct {
	svc *Service
}

func NewHandler(db models.DBTX) *Handler {
	repo := NewRepository(db)
	return &Handler{svc: NewService(repo)}
}

func (h *Handler) RegisterAssessmentRoutes(r chi.Router) {
	r.Post("/", h.CreateAssessment)
	r.Get("/{id}", h.GetAssessment)
	r.Get("/", h.ListAssessments)
	r.Patch("/{id}", h.UpdateAssessment)
	r.Delete("/{id}", h.DeleteAssessment)
}

func (h *Handler) CreateAssessment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	var req CreateAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	assessment, err := h.svc.CreateAssessment(ctx, req, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create assessment")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(assessment)
}

func (h *Handler) GetAssessment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	assessment, err := h.svc.GetAssessment(ctx, id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(assessment)
}

func (h *Handler) ListAssessments(w http.ResponseWriter, r *http.Request) {
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
	assessments, err := h.svc.ListAssessments(ctx, int32(countyID), int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(assessments)
}

func (h *Handler) UpdateAssessment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateAssessmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	assessment, err := h.svc.UpdateAssessment(ctx, id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(assessment)
}

func (h *Handler) DeleteAssessment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	if err := h.svc.DeleteAssessment(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Assessment Items Handlers
func (h *Handler) CreateAssessmentItem(w http.ResponseWriter, r *http.Request) {
	assessmentID := chi.URLParam(r, "id")
	var req CreateAssessmentItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.AssessmentID = assessmentID
	ctx := r.Context()
	item, err := h.svc.CreateAssessmentItem(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create assessment item")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) ListAssessmentItems(w http.ResponseWriter, r *http.Request) {
	assessmentID := chi.URLParam(r, "id")
	ctx := r.Context()
	items, err := h.svc.ListAssessmentItems(ctx, assessmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func (h *Handler) DeleteAssessmentItem(w http.ResponseWriter, r *http.Request) {
	assessmentID := chi.URLParam(r, "id")
	itemID := chi.URLParam(r, "item_id")
	ctx := r.Context()
	if err := h.svc.DeleteAssessmentItem(ctx, assessmentID, itemID); err != nil {
		if err.Error() == "assessment item not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Error().Err(err).Msg("Failed to delete assessment item")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}