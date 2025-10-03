package user

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sangkips/revenue-system/internal/db"
	"github.com/sangkips/revenue-system/internal/domain/user/models"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateUser(t *testing.T) {
	conn := db.SetUpTestDB(t)
	defer db.TeardownTestDB(t, conn)

	repo := NewRepository(conn)

	ctx := context.Background()

	params := models.InsertUserParams{
		CountyID:     sql.NullInt32{Valid: true, Int32: 1},
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FirstName:    "Test",
		LastName:     "User",
		Role:         "county_admin",
		IsActive:     sql.NullBool{Valid: true, Bool: true},
	}

	_, err := repo.CreateUser(ctx, params)
	assert.NoError(t, err)
	user, err := repo.GetUserByEmail(ctx, "test@example.com")
	assert.Equal(t, params.Email, user.Email)
}
