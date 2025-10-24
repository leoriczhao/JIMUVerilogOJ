package seed

import (
	"log"
	"verilog-oj/backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDefaultAdmin 创建默认管理员账户
// 仅在开发环境使用，用于 E2E 测试
func SeedDefaultAdmin(db *gorm.DB) error {
	// 检查是否已存在管理员用户
	var count int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	// 如果已有管理员，跳过
	if count > 0 {
		log.Println("Admin user already exists, skipping seed")
		return nil
	}

	// 创建默认管理员账户（用于开发和测试）
	log.Println("Creating default admin user for development/testing...")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Username: "admin",
		Email:    "admin@verilogoj.local",
		Password: string(hashedPassword),
		Nickname: "系统管理员",
		Role:     "admin",
		IsActive: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Printf("✅ Default admin user created (username: admin, password: admin123)")
	log.Printf("⚠️  WARNING: This is for development only! Change password in production!")
	return nil
}

// SeedCustomAdmin 创建自定义管理员账户
// 用于生产环境首次部署
func SeedCustomAdmin(db *gorm.DB, username, email, password string) error {
	// 检查是否已存在管理员用户
	var count int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin user already exists, skipping custom admin creation")
		return nil
	}

	if username == "" || email == "" || password == "" {
		log.Println("No custom admin credentials provided, skipping")
		return nil
	}

	log.Printf("Creating custom admin user: %s", username)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Nickname: "管理员",
		Role:     "admin",
		IsActive: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Printf("✅ Custom admin user created: %s", username)
	return nil
}
