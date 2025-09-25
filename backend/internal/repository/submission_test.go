package repository

import (
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/services"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSubmissionTestDB(t *testing.T) (*gorm.DB, services.UserRepository, services.ProblemRepository) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to memory db: %v", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Problem{}, &models.Submission{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}
	userRepo := NewUserRepository(db)
	problemRepo := NewProblemRepository(db)
	return db, userRepo, problemRepo
}

func TestSubmissionRepository_CreateAndGet(t *testing.T) {
	db, userRepo, problemRepo := setupSubmissionTestDB(t)
	repo := NewSubmissionRepository(db)

	user := &domain.User{Username: "submitter", Email: "submitter@test.com", Password: "pw"}
	userRepo.Create(user)
	problem := &domain.Problem{Title: "Submittable Problem"}
	problemRepo.Create(problem)

	submission := &domain.Submission{
		UserID:    user.ID,
		ProblemID: problem.ID,
		Code:      "module main; endmodule",
		Language:  "verilog",
	}
	err := repo.Create(submission)
	assert.NoError(t, err)
	assert.NotZero(t, submission.ID)

	retrieved, err := repo.GetByID(submission.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "verilog", retrieved.Language)
}

func TestSubmissionRepository_UpdateStatus(t *testing.T) {
	db, userRepo, problemRepo := setupSubmissionTestDB(t)
	repo := NewSubmissionRepository(db)
	user := &domain.User{Username: "submitter", Email: "submitter@test.com", Password: "pw"}
	userRepo.Create(user)
	problem := &domain.Problem{Title: "Submittable Problem"}
	problemRepo.Create(problem)
	submission := &domain.Submission{UserID: user.ID, ProblemID: problem.ID, Code: "code"}
	repo.Create(submission)

	err := repo.UpdateStatus(submission.ID, "Accepted", 100, 50, 1024, "", 10, 10)
	assert.NoError(t, err)

	retrieved, _ := repo.GetByID(submission.ID)
	assert.Equal(t, "Accepted", retrieved.Status)
	assert.Equal(t, 100, retrieved.Score)
}

func TestSubmissionRepository_List(t *testing.T) {
	db, userRepo, problemRepo := setupSubmissionTestDB(t)
	repo := NewSubmissionRepository(db)
	user1 := &domain.User{Username: "u1", Email: "u1@test.com", Password: "pw"}
	userRepo.Create(user1)
	user2 := &domain.User{Username: "u2", Email: "u2@test.com", Password: "pw"}
	userRepo.Create(user2)
	p1 := &domain.Problem{Title: "P1"}
	problemRepo.Create(p1)

	repo.Create(&domain.Submission{UserID: user1.ID, ProblemID: p1.ID, Status: "Accepted"})
	repo.Create(&domain.Submission{UserID: user2.ID, ProblemID: p1.ID, Status: "Wrong Answer"})
	repo.Create(&domain.Submission{UserID: user1.ID, ProblemID: p1.ID, Status: "Accepted"})

	// List all for user1
	list, total, err := repo.List(1, 10, user1.ID, 0, "")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)

	// List accepted for user1
	list, total, err = repo.List(1, 10, user1.ID, 0, "Accepted")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
}

func TestSubmissionRepository_CountAcceptedByUser(t *testing.T) {
	db, userRepo, problemRepo := setupSubmissionTestDB(t)
	repo := NewSubmissionRepository(db)
	user := &domain.User{Username: "u1", Email: "u1@test.com", Password: "pw"}
	userRepo.Create(user)
	p1 := &domain.Problem{Title: "P1"}
	problemRepo.Create(p1)
	p2 := &domain.Problem{Title: "P2"}
	problemRepo.Create(p2)

	repo.Create(&domain.Submission{UserID: user.ID, ProblemID: p1.ID, Status: "Accepted"})
	repo.Create(&domain.Submission{UserID: user.ID, ProblemID: p1.ID, Status: "Wrong Answer"})
	repo.Create(&domain.Submission{UserID: user.ID, ProblemID: p2.ID, Status: "Accepted"})

	count, err := repo.CountAcceptedByUser(user.ID, p1.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestSubmissionRepository_GetStats(t *testing.T) {
	db, userRepo, problemRepo := setupSubmissionTestDB(t)
	repo := NewSubmissionRepository(db)
	user := &domain.User{Username: "u1", Email: "u1@test.com", Password: "pw"}
	userRepo.Create(user)
	p1 := &domain.Problem{Title: "P1"}
	problemRepo.Create(p1)
	p2 := &domain.Problem{Title: "P2"}
	problemRepo.Create(p2)

	repo.Create(&domain.Submission{UserID: user.ID, ProblemID: p1.ID, Status: "Accepted"})
	repo.Create(&domain.Submission{UserID: user.ID, ProblemID: p1.ID, Status: "Wrong Answer"})
	repo.Create(&domain.Submission{UserID: user.ID, ProblemID: p2.ID, Status: "Accepted"})

	stats, err := repo.GetStats(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), stats["total_submissions"])
	assert.Equal(t, int64(2), stats["accepted_submissions"])
	assert.Equal(t, int64(2), stats["solved_problems"])
}