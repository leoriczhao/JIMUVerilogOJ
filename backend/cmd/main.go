package main

import (
	"fmt"
	"log"
	"verilog-oj/backend/internal"
	"verilog-oj/backend/internal/config"
	"verilog-oj/backend/internal/middleware"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/seed"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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

	// 初始化管理员账户
	if err := initializeAdmin(db, cfg); err != nil {
		log.Printf("Warning: Failed to initialize admin: %v", err)
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

	// 添加现有Service到中间件上下文
	r.Use(func(c *gin.Context) {
		c.Set("problem_service", app.Services.ProblemService)
		c.Set("forum_service", app.Services.ForumService)
		c.Next()
	})

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
		admin.Use(middleware.AuthRequired(), middleware.RequirePermission(middleware.PermManageSystem))
		{
			admin.GET("/whoami", app.Handlers.AdminHandler.WhoAmI)
			admin.GET("/stats", app.Handlers.AdminHandler.Stats)
			admin.PUT("/users/:id/role", app.Handlers.AdminHandler.UpdateUserRole)
		}

		// 用户相关路由
		users := v1.Group("/users")
		{
			users.POST("/register", app.Handlers.UserHandler.Register)
			users.POST("/login", app.Handlers.UserHandler.Login)
			// 获取个人信息：需要 user.profile.read 权限
			users.GET("/profile",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermUserProfileRead),
				app.Handlers.UserHandler.GetProfile)
			// 更新个人信息：需要 user.profile.update 权限
			users.PUT("/profile",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermUserProfileUpdate),
				app.Handlers.UserHandler.UpdateProfile)
			// 修改密码：需要 user.password.change 权限
			users.PUT("/password",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermUserPasswordChange),
				app.Handlers.UserHandler.ChangePassword)
		}

		// 题目相关路由
		problems := v1.Group("/problems")
		{
			// 获取题目列表：需要 problem.list 权限
			problems.GET("",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermProblemList),
				app.Handlers.ProblemHandler.ListProblems)
			// 获取题目详情：需要 problem.read 权限
			problems.GET("/:id",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermProblemRead),
				app.Handlers.ProblemHandler.GetProblem)
			// 创建题目：需要 problem.create 权限
			problems.POST("",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermProblemCreate),
				app.Handlers.ProblemHandler.CreateProblem)
			// 更新题目：需要 problem.update.own 权限或管理员角色
			problems.PUT("/:id",
				middleware.AuthRequired(),
				middleware.RequireOwnershipOrPermission(
					middleware.PermProblemUpdateAll,
					middleware.GetProblemOwner("id"),
				),
				app.Handlers.ProblemHandler.UpdateProblem)
			// 删除题目：需要 problem.delete.own 权限或管理员角色
			problems.DELETE("/:id",
				middleware.AuthRequired(),
				middleware.RequireOwnershipOrPermission(
					middleware.PermProblemDeleteAll,
					middleware.GetProblemOwner("id"),
				),
				app.Handlers.ProblemHandler.DeleteProblem)

			// 测试用例相关路由
			// 查看测试用例：支持可选认证（用户看样例，教师和管理员看全部）
			problems.GET("/:id/testcases",
				middleware.OptionalAuth(),
				app.Handlers.ProblemHandler.GetTestCases)
			// 添加测试用例：需要 testcase.create 权限或管理员角色
			problems.POST("/:id/testcases",
				middleware.AuthRequired(),
				middleware.RequireOwnershipOrPermission(
					middleware.PermTestcaseCreate,
					middleware.GetProblemOwner("id"),
				),
				app.Handlers.ProblemHandler.AddTestCase)

			// 题目提交记录：需要 submission.list 权限
			problems.GET("/:id/submissions",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermSubmissionList),
				app.Handlers.SubmissionHandler.GetProblemSubmissions)
		}

		// 提交相关路由
		submissions := v1.Group("/submissions")
		{
			// 获取公开提交列表：支持可选认证，匿名用户看基础信息
			submissions.GET("",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermSubmissionList),
				app.Handlers.SubmissionHandler.ListSubmissions)
			// 获取提交详情：需要 submission.read 权限，支持可选认证
			submissions.GET("/:id",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermSubmissionRead),
				app.Handlers.SubmissionHandler.GetSubmission)
			// 创建提交：需要 submission.create 权限
			submissions.POST("",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermSubmissionCreate),
				app.Handlers.SubmissionHandler.CreateSubmission)
			// 删除提交：需要 submission.delete 权限（仅管理员）
			submissions.DELETE("/:id",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermSubmissionDelete),
				app.Handlers.SubmissionHandler.DeleteSubmission)

			// 获取当前用户提交记录：需要 submission.list 权限
			submissions.GET("/user",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermSubmissionList),
				app.Handlers.SubmissionHandler.GetUserSubmissions)
			// 获取提交统计：需要 submission.read 权限
			submissions.GET("/stats",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermSubmissionRead),
				app.Handlers.SubmissionHandler.GetSubmissionStats)
		}

		// 论坛相关路由
		forum := v1.Group("/forum")
		{
			// 获取帖子列表：支持可选认证，匿名用户可查看公开帖子
			forum.GET("/posts",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermForumPostRead),
				app.Handlers.ForumHandler.ListPosts)
			// 获取帖子详情：需要 forum.post.read 权限，支持可选认证
			forum.GET("/posts/:id",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermForumPostRead),
				app.Handlers.ForumHandler.GetPost)
			// 创建帖子：需要 forum.post.create 权限
			forum.POST("/posts",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermForumPostCreate),
				app.Handlers.ForumHandler.CreatePost)
			// 更新帖子：需要 forum.edit.own 权限或管理员角色
			forum.PUT("/posts/:id",
				middleware.AuthRequired(),
				middleware.RequireOwnershipOrPermission(
					middleware.PermForumEditAll,
					middleware.GetForumPostOwner("id"),
				),
				app.Handlers.ForumHandler.UpdatePost)
			// 删除帖子：需要 forum.edit.own 权限或管理员角色
			forum.DELETE("/posts/:id",
				middleware.AuthRequired(),
				middleware.RequireOwnershipOrPermission(
					middleware.PermForumDelete,
					middleware.GetForumPostOwner("id"),
				),
				app.Handlers.ForumHandler.DeletePost)

			// 获取帖子回复列表：支持可选认证，匿名用户可查看公开回复
			forum.GET("/posts/:id/replies",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermForumReplyRead),
				app.Handlers.ForumHandler.ListReplies)
			// 创建回复：需要 forum.reply.create 权限
			forum.POST("/posts/:id/replies",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermForumReplyCreate),
				app.Handlers.ForumHandler.CreateReply)
		}

		// 新闻相关路由
		news := v1.Group("/news")
		{
			// 获取新闻列表：需要 news.list 权限
			news.GET("",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermNewsList),
				app.Handlers.NewsHandler.ListNews)
			// 获取新闻详情：需要 news.read 权限
			news.GET("/:id",
				middleware.OptionalAuth(),
				middleware.OptionalAuthPermission(middleware.PermNewsRead),
				app.Handlers.NewsHandler.GetNews)
			// 创建新闻：需要 news.create 权限
			news.POST("",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermNewsCreate),
				app.Handlers.NewsHandler.CreateNews)
			// 更新新闻：需要 news.update 权限
			news.PUT("/:id",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermNewsUpdate),
				app.Handlers.NewsHandler.UpdateNews)
			// 删除新闻：需要 news.delete 权限
			news.DELETE("/:id",
				middleware.AuthRequired(),
				middleware.RequirePermission(middleware.PermNewsDelete),
				app.Handlers.NewsHandler.DeleteNews)
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info), // Enable detailed SQL logging
		DisableForeignKeyConstraintWhenMigrating: true,                                // Disable foreign key constraints
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

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

// initializeAdmin 初始化管理员账户
func initializeAdmin(db *gorm.DB, cfg *config.Config) error {
	// 开发模式：创建默认管理员 (admin/admin123) 用于 E2E 测试
	if cfg.Server.Mode == "debug" {
		log.Println("Development mode: Creating default admin user for testing...")
		return seed.SeedDefaultAdmin(db)
	}

	// 生产模式：如果配置了自定义管理员，则创建
	if cfg.InitAdmin.Username != "" {
		log.Println("Production mode: Creating custom admin user...")
		return seed.SeedCustomAdmin(
			db,
			cfg.InitAdmin.Username,
			cfg.InitAdmin.Email,
			cfg.InitAdmin.Password,
		)
	}

	log.Println("No admin initialization required")
	return nil
}
