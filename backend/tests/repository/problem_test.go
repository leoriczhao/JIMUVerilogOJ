package repository_test

import (
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/repository"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupProblemTestDB creates a new in-memory SQLite database for testing the problem repository.
func setupProblemTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to memory db: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.Problem{}, &models.TestCase{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}
	return db
}

func TestProblemRepository_CreateAndGet(t *testing.T) {
	db := setupProblemTestDB(t)
	repo := repository.NewProblemRepository(db)

	problem := &domain.Problem{
		Title:       "Test Problem",
		Description: "Problem Description",
		Difficulty:  "Easy",
		Category:    "Array",
	}
	err := repo.Create(problem)
	assert.NoError(t, err)
	assert.NotZero(t, problem.ID)

	retrieved, err := repo.GetByID(problem.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "Test Problem", retrieved.Title)
}

func TestProblemRepository_Update(t *testing.T) {
	db := setupProblemTestDB(t)
	repo := repository.NewProblemRepository(db)

	problem := &domain.Problem{Title: "Original", Description: "Original"}
	repo.Create(problem)

	problem.Title = "Updated"
	err := repo.Update(problem)
	assert.NoError(t, err)

	retrieved, _ := repo.GetByID(problem.ID)
	assert.Equal(t, "Updated", retrieved.Title)
}

func TestProblemRepository_Delete(t *testing.T) {
	db := setupProblemTestDB(t)
	repo := repository.NewProblemRepository(db)

	problem := &domain.Problem{Title: "To Be Deleted"}
	repo.Create(problem)

	err := repo.Delete(problem.ID)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(problem.ID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestProblemRepository_List(t *testing.T) {
	db := setupProblemTestDB(t)
	repo := repository.NewProblemRepository(db)

	repo.Create(&domain.Problem{Title: "P1", Difficulty: "Easy", Category: "Array"})
	repo.Create(&domain.Problem{Title: "P2", Difficulty: "Medium", Category: "String"})
	repo.Create(&domain.Problem{Title: "Another Easy One", Difficulty: "Easy", Category: "Array"})

	// List all
	list, total, err := repo.List(1, 10, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list, 3)

	// Filter by difficulty
	list, total, err = repo.List(1, 10, map[string]interface{}{"difficulty": "Easy"})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)
}

func TestProblemRepository_TestCases(t *testing.T) {
	db := setupProblemTestDB(t)
	repo := repository.NewProblemRepository(db)

	problem := &domain.Problem{Title: "Problem With Cases"}
	repo.Create(problem)

	// Create test cases
	tc1 := &domain.TestCase{ProblemID: problem.ID, Input: "in1", Output: "out1"}
	repo.CreateTestCase(tc1)
	tc2 := &domain.TestCase{ProblemID: problem.ID, Input: "in2", Output: "out2"}
	repo.CreateTestCase(tc2)

	// Get test cases
	cases, err := repo.GetTestCases(problem.ID)
	assert.NoError(t, err)
	assert.Len(t, cases, 2)
	assert.Equal(t, "in1", cases[0].Input)

	// Delete test cases
	err = repo.DeleteTestCases(problem.ID)
	assert.NoError(t, err)
	cases, err = repo.GetTestCases(problem.ID)
	assert.NoError(t, err)
	assert.Len(t, cases, 0)
}

func TestProblemRepository_UpdateCounts(t *testing.T) {
	db := setupProblemTestDB(t)
	repo := repository.NewProblemRepository(db)

	problem := &domain.Problem{Title: "Counter Problem"}
	repo.Create(problem)

	// Update submit count
	repo.UpdateSubmitCount(problem.ID, 1)
	repo.UpdateSubmitCount(problem.ID, 1)
	retrieved, _ := repo.GetByID(problem.ID)
	assert.Equal(t, 2, retrieved.SubmitCount)

	// Update accepted count
	repo.UpdateAcceptedCount(problem.ID, 1)
	retrieved, _ = repo.GetByID(problem.ID)
	assert.Equal(t, 1, retrieved.AcceptedCount)
}
