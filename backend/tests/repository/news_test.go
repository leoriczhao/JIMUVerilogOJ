package repository_test

import (
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/repository"
	"verilog-oj/backend/internal/services"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupNewsTestDB creates a new in-memory SQLite database for testing the news repository.
func setupNewsTestDB(t *testing.T) (*gorm.DB, services.UserRepository) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to memory db: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.News{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	return db, userRepo
}

func TestNewsRepository_CreateAndGet(t *testing.T) {
	db, userRepo := setupNewsTestDB(t)
	repo := repository.NewNewsRepository(db)

	user := &domain.User{Username: "news_author", Email: "news@test.com", Password: "pw"}
	userRepo.Create(user)

	news := &domain.News{
		Title:       "Test News",
		Content:     "News Content",
		AuthorID:    user.ID,
		IsFeatured:  true,
	}

	err := repo.Create(news)
	assert.NoError(t, err)
	assert.NotZero(t, news.ID)

	retrieved, err := repo.GetByID(news.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "Test News", retrieved.Title)
	assert.True(t, retrieved.IsFeatured)
}

func TestNewsRepository_Update(t *testing.T) {
	db, userRepo := setupNewsTestDB(t)
	repo := repository.NewNewsRepository(db)

	user := &domain.User{Username: "news_author", Email: "news@test.com", Password: "pw"}
	userRepo.Create(user)

	news := &domain.News{Title: "Original", Content: "Original", AuthorID: user.ID}
	repo.Create(news)

	news.Title = "Updated"
	news.IsFeatured = true
	err := repo.Update(news)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(news.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", retrieved.Title)
	assert.True(t, retrieved.IsFeatured)
}

func TestNewsRepository_Delete(t *testing.T) {
	db, userRepo := setupNewsTestDB(t)
	repo := repository.NewNewsRepository(db)

	user := &domain.User{Username: "news_author", Email: "news@test.com", Password: "pw"}
	userRepo.Create(user)
	news := &domain.News{Title: "To Be Deleted", AuthorID: user.ID}
	repo.Create(news)

	err := repo.Delete(news.ID)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(news.ID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestNewsRepository_List(t *testing.T) {
	db, userRepo := setupNewsTestDB(t)
	repo := repository.NewNewsRepository(db)

	user := &domain.User{Username: "news_author", Email: "news@test.com", Password: "pw"}
	userRepo.Create(user)

	repo.Create(&domain.News{Title: "N1", AuthorID: user.ID, IsFeatured: true})
	repo.Create(&domain.News{Title: "N2", AuthorID: user.ID})
	repo.Create(&domain.News{Title: "N3", AuthorID: user.ID})

	list, total, err := repo.List(1, 10, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list, 3)
	// Check ordering
	assert.Equal(t, "N1", list[0].Title)
	assert.True(t, list[0].IsFeatured)
}

func TestNewsRepository_IncrementViewCount(t *testing.T) {
	db, userRepo := setupNewsTestDB(t)
	repo := repository.NewNewsRepository(db)

	user := &domain.User{Username: "news_author", Email: "news@test.com", Password: "pw"}
	userRepo.Create(user)
	news := &domain.News{Title: "Views", AuthorID: user.ID}
	repo.Create(news)

	err := repo.IncrementViewCount(news.ID)
	assert.NoError(t, err)

	retrieved, err := repo.GetByID(news.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, retrieved.ViewCount)
}