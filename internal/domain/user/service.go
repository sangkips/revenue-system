package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sangkips/revenue-system/internal/domain/user/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (models.User, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return models.User{}, errors.New("required fields missing")
	}

	validRoles := map[string]bool{
		"super_admin":     true,
		"county_admin":    true,
		"department_head": true,
		"collector":       true,
		"auditor":         true,
	}
	if !validRoles[req.Role] {
		return models.User{}, errors.New("invalid role")
	}

	if req.Role != "super_admin" && req.CountyID == nil {
		return models.User{}, errors.New("county_id is required for non-super-admin roles")
	}

	hashedPsswd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	var countyID sql.NullInt32
	if req.CountyID != nil {
		countyID = sql.NullInt32{Int32: *req.CountyID, Valid: true}
	}

	params := models.InsertUserParams{
		CountyID:     countyID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPsswd),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PhoneNumber:  sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		Role:         req.Role,
		EmployeeID:   sql.NullString{String: req.EmployeeID, Valid: req.EmployeeID != ""},
		Department:   sql.NullString{String: req.Department, Valid: req.Department != ""},
		IsActive:     sql.NullBool{Bool: true, Valid: true},
	}

	return s.repo.CreateUser(ctx, params)
}

func (s *Service) ListUsers(ctx context.Context, userRole string, userCountyID *int32, limit, offset int32) (interface{}, error) {

	if userRole == "super_admin" {
		return s.repo.ListAllUsers(ctx, models.ListAllUsersParams{Limit: limit, Offset: offset})
	}

	if userCountyID == nil {
		return nil, errors.New("county_id is required for non-super-admin roles")
	}

	return s.repo.ListUsers(ctx, models.ListUsersParams{
		CountyID: sql.NullInt32{Int32: *userCountyID, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
}

func (s *Service) GetUser(ctx context.Context, id string) (models.GetUserByIDRow, error) {
	return s.repo.GetUserByID(ctx, id)
}

type CreateUserRequest struct {
	CountyID    *int32 `json:"county_id,omitempty"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	EmployeeID  string `json:"employee_id"`
	Department  string `json:"department"`
}

type UpdateUserRequest struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	EmployeeID  string `json:"employee_id"`
	Department  string `json:"department"`
	IsActive    bool   `json:"is_active"`
}
