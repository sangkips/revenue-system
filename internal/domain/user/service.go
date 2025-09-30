package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
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
	if req.Email == "" || req.Password == "" || req.Role == "" {
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

func (s *Service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	params := models.UpdateUserParams{
		ID: userID,
	}

	if req.Email != nil {
		params.UpdateEmail = true
		params.Email = *req.Email
	}

	if req.FirstName != nil {
		params.UpdateFirstName = true
		params.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		params.UpdateLastName = true
		params.LastName = *req.LastName
	}

	if req.PhoneNumber != nil {
		params.UpdatePhoneNumber = true
		params.PhoneNumber = sql.NullString{String: *req.PhoneNumber, Valid: *req.PhoneNumber != ""}
	}

	if req.Role != nil {
		params.UpdateRole = true
		params.Role = *req.Role
	}

	if req.EmployeeID != nil {
		params.UpdateEmployeeID = true
		params.EmployeeID = sql.NullString{String: *req.EmployeeID, Valid: *req.EmployeeID != ""}
	}

	if req.Department != nil {
		params.UpdateDepartment = true
		params.Department = sql.NullString{String: *req.Department, Valid: *req.Department != ""}
	}

	if req.IsActive != nil {
		params.UpdateIsActive = true
		params.IsActive = sql.NullBool{Bool: *req.IsActive, Valid: true}
	}

	return s.repo.UpdateUser(ctx, params)
}

func(s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

type CreateUserRequest struct {
	CountyID    *int32 `json:"county_id,omitempty"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	EmployeeID  string `json:"employee_id"`
	Department  string `json:"department"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

type UpdateUserRequest struct {
	Email       *string `json:"email,omitempty"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Role        *string `json:"role,omitempty"`
	EmployeeID  *string `json:"employee_id,omitempty"`
	Department  *string `json:"department,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}
