package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/sangkips/revenue-system/internal/domain/user/models"
	"github.com/sangkips/revenue-system/internal/middleware/auth"
)

type Handler struct {
	svc *Service
}

type UserResponse struct {
	ID          string  `json:"id"`
	CountyID    *int32  `json:"county_id"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	PhoneNumber *string `json:"phone_number"`
	Role        string  `json:"role"`
	EmployeeID  *string `json:"employee_id"`
	Department  *string `json:"department"`
	IsActive    *bool   `json:"is_active"`
	LastLogin   *string `json:"last_login"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

func NewHandler(db models.DBTX) *Handler {
	repo := NewRepository(db)

	return &Handler{svc: NewService(repo)}
}

func (h Handler) RegisterUserRoutes(r chi.Router) {
	r.Post("/", h.CreateUser)
	r.Get("/{id}", h.GetUser)
	r.Get("/", h.ListUsers)
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

func (h Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch user")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := convertGetUserByIDRowToResponse(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (h Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	if limit == 0 {
		limit = 30
	}

	ctx := r.Context()

	userRole, ok := ctx.Value(auth.UserRoleKey).(string)
	if !ok {
		http.Error(w, "user role not found in context", http.StatusUnauthorized)
		return
	}

	var userCountyID *int32
	if countyIDValue := ctx.Value(auth.UserCountyIDKey); countyIDValue != nil {
		if countyID, ok := countyIDValue.(int32); ok {
			userCountyID = &countyID
		}
	}

	users, err := h.svc.ListUsers(ctx, userRole, userCountyID, int32(limit), int32(offset))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert database models to clean JSON responses
	var responses []UserResponse
	switch v := users.(type) {
	case []models.ListUsersRow:
		for _, user := range v {
			responses = append(responses, convertListUsersRowToResponse(user))
		}
	case []models.ListAllUsersRow:
		for _, user := range v {
			responses = append(responses, convertListAllUsersRowToResponse(user))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}

// Convert ListUsersRow to UserResponse
func convertListUsersRowToResponse(user models.ListUsersRow) UserResponse {
	response := UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	// Handle nullable fields
	if user.CountyID.Valid {
		response.CountyID = &user.CountyID.Int32
	}
	if user.PhoneNumber.Valid {
		response.PhoneNumber = &user.PhoneNumber.String
	}
	if user.EmployeeID.Valid {
		response.EmployeeID = &user.EmployeeID.String
	}
	if user.Department.Valid {
		response.Department = &user.Department.String
	}
	if user.IsActive.Valid {
		response.IsActive = &user.IsActive.Bool
	}
	if user.LastLogin.Valid {
		lastLogin := user.LastLogin.Time.Format("2006-01-02T15:04:05Z")
		response.LastLogin = &lastLogin
	}
	if user.CreatedAt.Valid {
		createdAt := user.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
		response.CreatedAt = &createdAt
	}
	if user.UpdatedAt.Valid {
		updatedAt := user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
		response.UpdatedAt = &updatedAt
	}

	return response
}

// Convert GetUserByIDRow to UserResponse
func convertGetUserByIDRowToResponse(user models.GetUserByIDRow) UserResponse {
	response := UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	// Handle nullable fields
	if user.CountyID.Valid {
		response.CountyID = &user.CountyID.Int32
	}
	if user.PhoneNumber.Valid {
		response.PhoneNumber = &user.PhoneNumber.String
	}
	if user.EmployeeID.Valid {
		response.EmployeeID = &user.EmployeeID.String
	}
	if user.Department.Valid {
		response.Department = &user.Department.String
	}
	if user.IsActive.Valid {
		response.IsActive = &user.IsActive.Bool
	}
	if user.LastLogin.Valid {
		lastLogin := user.LastLogin.Time.Format("2006-01-02T15:04:05Z")
		response.LastLogin = &lastLogin
	}
	if user.CreatedAt.Valid {
		createdAt := user.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
		response.CreatedAt = &createdAt
	}
	if user.UpdatedAt.Valid {
		updatedAt := user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
		response.UpdatedAt = &updatedAt
	}

	return response
}

// Convert ListAllUsersRow to UserResponse
func convertListAllUsersRowToResponse(user models.ListAllUsersRow) UserResponse {
	response := UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}

	// Handle nullable fields
	if user.CountyID.Valid {
		response.CountyID = &user.CountyID.Int32
	}
	if user.PhoneNumber.Valid {
		response.PhoneNumber = &user.PhoneNumber.String
	}
	if user.EmployeeID.Valid {
		response.EmployeeID = &user.EmployeeID.String
	}
	if user.Department.Valid {
		response.Department = &user.Department.String
	}
	if user.IsActive.Valid {
		response.IsActive = &user.IsActive.Bool
	}
	if user.LastLogin.Valid {
		lastLogin := user.LastLogin.Time.Format("2006-01-02T15:04:05Z")
		response.LastLogin = &lastLogin
	}
	if user.CreatedAt.Valid {
		createdAt := user.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
		response.CreatedAt = &createdAt
	}
	if user.UpdatedAt.Valid {
		updatedAt := user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
		response.UpdatedAt = &updatedAt
	}

	return response
}
