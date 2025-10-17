package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sangkips/revenue-system/internal/domain/user/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      Repository
	secretKey []byte
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.InsertUserParams) (models.User, error)
}

func NewAuthService(repo Repository, secretKey string) *AuthService {
	return &AuthService{repo: repo, secretKey: []byte(secretKey)}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (s AuthService) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return LoginResponse{}, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResponse{}, errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	if user.CountyID.Valid {
		claims["county_id"] = user.CountyID.Int32
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{Token: tokenString}, nil
}

type RegisterRequest struct {
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

func (s AuthService) Register(ctx context.Context, req RegisterRequest) (models.User, error) {
	if req.Role == "" {
		req.Role = "user"
	}
	
	if req.Email == "" || req.Password == "" {
		return models.User{}, errors.New("required fields missing")
	}

	// Validate role
	validRoles := map[string]bool{
		"user":			   true,
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
	user, err := s.repo.CreateUser(ctx, params)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
