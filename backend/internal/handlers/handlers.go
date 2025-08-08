package handlers

import (
	"verilog-oj/backend/internal/services"
)

// Handlers 包含所有Handler的结构体
type Handlers struct {
	UserHandler       *UserHandler
	ProblemHandler    *ProblemHandler
	SubmissionHandler *SubmissionHandler
	ForumHandler      *ForumHandler
	NewsHandler       *NewsHandler
	AdminHandler      *AdminHandler
}

// NewHandlers 创建Handlers实例
func NewHandlers(
	userService *services.UserService,
	problemService *services.ProblemService,
	submissionService *services.SubmissionService,
	forumService *services.ForumService,
	newsService *services.NewsService,
) *Handlers {
	return &Handlers{
		UserHandler:       NewUserHandler(userService),
		ProblemHandler:    NewProblemHandler(problemService),
		SubmissionHandler: NewSubmissionHandler(submissionService),
		ForumHandler:      NewForumHandler(forumService),
		NewsHandler:       NewNewsHandler(newsService),
		AdminHandler:      NewAdminHandler(),
	}
}
