package main

import (
	"fmt"
	"log"
	"verilog-oj/backend/internal"
	"verilog-oj/backend/internal/config"
	"verilog-oj/backend/internal/middleware"
	"verilog-oj/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @title Verilog OJ API
// @version 1.0
// @description Verilog在线判题系统API文档
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 自动迁移数据库表
	if errMigrate := autoMigrate(db); errMigrate != nil {
		log.Fatal("Failed to migrate database:", errMigrate)
	}

	// 使用 wire 初始化应用
	app, err := internal.InitializeApp(db)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	r := gin.New()

	// 添加中间件
	r.Use(middleware.CORS())
	r.Use(middleware.APILogger()) // 使用文件日志中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.APIRateLimit()) // 添加API速率限制

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Verilog OJ API is running",
		})
	})

	// API路由组
	v1 := r.Group("/api/v1")
	{
		// 管理端路由（需要管理员权限）
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired(), middleware.AdminOnly())
		{
			admin.GET("/whoami", app.Handlers.AdminHandler.WhoAmI)
			admin.GET("/stats", app.Handlers.AdminHandler.Stats)
		}

		// 用户相关路由
		users := v1.Group("/users")
		{
			users.POST("/register", app.Handlers.UserHandler.Register)
			users.POST("/login", app.Handlers.UserHandler.Login)
			users.GET("/profile", middleware.AuthRequired(), app.Handlers.UserHandler.GetProfile)
			users.PUT("/profile", middleware.AuthRequired(), app.Handlers.UserHandler.UpdateProfile)
		}

		// 题目相关路由
		problems := v1.Group("/problems")
		{
			problems.GET("", app.Handlers.ProblemHandler.ListProblems)
			problems.GET("/:id", app.Handlers.ProblemHandler.GetProblem)
			problems.POST("", middleware.AuthRequired(), app.Handlers.ProblemHandler.CreateProblem)
			problems.PUT("/:id", middleware.AuthRequired(), app.Handlers.ProblemHandler.UpdateProblem)
			problems.DELETE("/:id", middleware.AuthRequired(), app.Handlers.ProblemHandler.DeleteProblem)

			// 测试用例相关路由
			problems.GET("/:id/testcases", app.Handlers.ProblemHandler.GetTestCases)
			problems.POST("/:id/testcases", middleware.AuthRequired(), app.Handlers.ProblemHandler.AddTestCase)

			// 题目提交记录路由
			problems.GET("/:id/submissions", app.Handlers.SubmissionHandler.GetProblemSubmissions)
		}

		// 提交相关路由
		submissions := v1.Group("/submissions")
		{
			submissions.GET("", app.Handlers.SubmissionHandler.ListSubmissions)
			submissions.GET("/:id", app.Handlers.SubmissionHandler.GetSubmission)
			submissions.POST("", middleware.AuthRequired(), app.Handlers.SubmissionHandler.CreateSubmission)
			submissions.DELETE("/:id", middleware.AuthRequired(), app.Handlers.SubmissionHandler.DeleteSubmission)

			// 用户提交记录
			submissions.GET("/user", middleware.AuthRequired(), app.Handlers.SubmissionHandler.GetUserSubmissions)
			// 提交统计
			submissions.GET("/stats", middleware.AuthRequired(), app.Handlers.SubmissionHandler.GetSubmissionStats)
		}

		// 论坛相关路由
		forum := v1.Group("/forum")
		{
			forum.GET("/posts", app.Handlers.ForumHandler.ListPosts)
			forum.GET("/posts/:id", app.Handlers.ForumHandler.GetPost)
			forum.POST("/posts", middleware.AuthRequired(), app.Handlers.ForumHandler.CreatePost)
			forum.PUT("/posts/:id", middleware.AuthRequired(), app.Handlers.ForumHandler.UpdatePost)
			forum.DELETE("/posts/:id", middleware.AuthRequired(), app.Handlers.ForumHandler.DeletePost)

			forum.GET("/posts/:id/replies", app.Handlers.ForumHandler.ListReplies)
			forum.POST("/posts/:id/replies", middleware.AuthRequired(), app.Handlers.ForumHandler.CreateReply)
		}

		// 新闻相关路由
		news := v1.Group("/news")
		{
			news.GET("", app.Handlers.NewsHandler.ListNews)
			news.GET("/:id", app.Handlers.NewsHandler.GetNews)
			news.POST("", middleware.AuthRequired(), app.Handlers.NewsHandler.CreateNews)
			news.PUT("/:id", middleware.AuthRequired(), app.Handlers.NewsHandler.UpdateNews)
			news.DELETE("/:id", middleware.AuthRequired(), app.Handlers.NewsHandler.DeleteNews)
		}
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDB 初始化数据库连接
func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Database,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.TestCase{},
		&models.Submission{},
		&models.ForumPost{},
		&models.ForumReply{},
		&models.ForumLike{},
		&models.News{},
	)
}
