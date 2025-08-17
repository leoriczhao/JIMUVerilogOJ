package services

// Services 包含所有服务
type Services struct {
	UserService       *UserService
	ProblemService    *ProblemService
	SubmissionService *SubmissionService
	ForumService      *ForumService
	NewsService       *NewsService
	AdminService      *AdminService
}

// NewServices 创建Services实例
func NewServices(
	userRepo UserRepository,
	problemRepo ProblemRepository,
	submissionRepo SubmissionRepository,
	forumRepo ForumRepository,
	newsRepo NewsRepository,
	adminRepo AdminRepository,
) *Services {
	return &Services{
		UserService:       NewUserService(userRepo),
		ProblemService:    NewProblemService(problemRepo),
		SubmissionService: NewSubmissionService(submissionRepo, problemRepo, userRepo),
		ForumService:      NewForumService(forumRepo, userRepo),
		NewsService:       NewNewsService(newsRepo, userRepo),
		AdminService:      NewAdminService(adminRepo),
	}
}
