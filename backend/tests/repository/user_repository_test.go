package repository

import (
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/repository"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建一个用于测试的内存SQLite数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("无法连接到内存数据库: %v", err)
	}

	// 自动迁移模式
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("无法迁移数据库: %v", err)
	}

	return db
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	// 1. 创建一个初始用户
	initialUser := &domain.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		Nickname:  "Initial Nickname",
		School:    "Initial School",
		StudentID: "Initial-001",
	}
	err := repo.Create(initialUser)
	assert.NoError(t, err)
	assert.NotZero(t, initialUser.ID)

	// 2. 更新用户信息
	updateInfo := &domain.User{
		ID:        initialUser.ID,
		Nickname:  "Updated Nickname",
		School:    "Updated School",
		StudentID: "Updated-002",
	}
	err = repo.Update(updateInfo)
	assert.NoError(t, err)

	// 3. 从数据库中检索用户以验证更新
	updatedUser, err := repo.GetByID(initialUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)

	// 4. 断言
	assert.Equal(t, "Updated Nickname", updatedUser.Nickname)
	assert.Equal(t, "Updated School", updatedUser.School)
	assert.Equal(t, "Updated-002", updatedUser.StudentID)
	// 这个断言会失败，因为 `Update` 方法中的 `db.Save` 会用空字符串覆盖现有的 email。
	assert.Equal(t, "test@example.com", updatedUser.Email, "Email不应该被清空")
}
