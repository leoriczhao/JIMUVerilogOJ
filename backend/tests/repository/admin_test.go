package repository_test

import (
	"testing"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/repository"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupAdminTestDB creates a new in-memory SQLite database for testing the admin repository.
func setupAdminTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to memory db: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Problem{}, &models.Submission{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	return db
}

func TestAdminRepositoryImpl_CountUsers(t *testing.T) {
	db := setupAdminTestDB(t)
	repo := repository.NewAdminRepository(db)

	// Case 1: No users
	count, err := repo.CountUsers()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Case 2: Some users
	db.Create(&models.User{Username: "user1", Email: "user1@test.com", Password: "password"})
	db.Create(&models.User{Username: "user2", Email: "user2@test.com", Password: "password"})

	count, err = repo.CountUsers()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestAdminRepositoryImpl_CountProblems(t *testing.T) {
	db := setupAdminTestDB(t)
	repo := repository.NewAdminRepository(db)

	// Case 1: No problems
	count, err := repo.CountProblems()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Case 2: Some problems
	db.Create(&models.Problem{Title: "Problem 1", Description: "Desc 1"})
	db.Create(&models.Problem{Title: "Problem 2", Description: "Desc 2"})
	db.Create(&models.Problem{Title: "Problem 3", Description: "Desc 3"})

	count, err = repo.CountProblems()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestAdminRepositoryImpl_CountSubmissions(t *testing.T) {
	db := setupAdminTestDB(t)
	repo := repository.NewAdminRepository(db)

	// Case 1: No submissions
	count, err := repo.CountSubmissions()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Case 2: Some submissions
	db.Create(&models.Submission{UserID: 1, ProblemID: 1, Code: "code1"})
	db.Create(&models.Submission{UserID: 1, ProblemID: 2, Code: "code2"})

	count, err = repo.CountSubmissions()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}