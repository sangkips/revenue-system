package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/sangkips/revenue-system/internal/domain/user/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateUser(t *testing.T) {
	repo := &MockRepository{}
	svc := NewService(repo)
	handler := &Handler{svc: svc}

	countyID := int32(1)
	reqBody := CreateUserRequest{
		CountyID: &countyID,
		Email:    "test@example.com",
		Password: "test123",
		Role:     "county_admin",
	}
	body, _ := json.Marshal(reqBody)
	repo.On("CreateUser", mock.Anything, mock.AnythingOfType("models.InsertUserParams")).Return(models.User{}, nil)

	r := chi.NewRouter()
	r.Post("/", handler.CreateUser)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	repo.AssertExpectations(t)
}
