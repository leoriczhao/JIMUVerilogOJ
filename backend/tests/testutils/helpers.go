package testutils

import (
	"fmt"
	"testing"
	"time"
	"verilog-oj/backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// CreateTestUser 创建测试用户
func CreateTestUser(t *testing.T, username, email string) *domain.User {
	t.Helper()
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	
	return &domain.User{
		ID:          1,
		Username:    username,
		Email:       email,
		Password:    string(hashedPassword),
		Nickname:    "Test User",
		Role:        "student",
		IsActive:    true,
		Rating:      1200,
		Solved:      0,
		Submitted:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}
}

// CreateTestUsers 创建多个测试用户
func CreateTestUsers(t *testing.T, count int) []*domain.User {
	t.Helper()
	
	users := make([]*domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = CreateTestUser(t, 
			fmt.Sprintf("testuser%d", i+1),
			fmt.Sprintf("test%d@example.com", i+1),
		)
		users[i].ID = uint(i + 1)
	}
	return users
}

// AssertUserEqual 断言两个用户对象相等
func AssertUserEqual(t *testing.T, expected, actual *domain.User) {
	t.Helper()
	
	if expected == nil && actual == nil {
		return
	}
	
	if expected == nil || actual == nil {
		t.Errorf("One of the users is nil: expected=%v, actual=%v", expected, actual)
		return
	}
	
	if expected.ID != actual.ID {
		t.Errorf("User ID mismatch: expected=%d, actual=%d", expected.ID, actual.ID)
	}
	
	if expected.Username != actual.Username {
		t.Errorf("Username mismatch: expected=%s, actual=%s", expected.Username, actual.Username)
	}
	
	if expected.Email != actual.Email {
		t.Errorf("Email mismatch: expected=%s, actual=%s", expected.Email, actual.Email)
	}
	
	if expected.Role != actual.Role {
		t.Errorf("Role mismatch: expected=%s, actual=%s", expected.Role, actual.Role)
	}
	
	if expected.IsActive != actual.IsActive {
		t.Errorf("IsActive mismatch: expected=%t, actual=%t", expected.IsActive, actual.IsActive)
	}
}

// GenerateHashedPassword 生成加密密码
func GenerateHashedPassword(t *testing.T, password string) string {
	t.Helper()
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

// TimeNow 返回当前时间（用于测试中的时间一致性）
func TimeNow() time.Time {
	return time.Now().Truncate(time.Second)
}

// TimeEqual 比较两个时间是否相等（忽略纳秒差异）
func TimeEqual(t1, t2 time.Time) bool {
	return t1.Truncate(time.Second).Equal(t2.Truncate(time.Second))
}