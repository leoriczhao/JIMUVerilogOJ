package dto

import (
	"encoding/json"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
)

// UserToResponse 将User模型转换为UserResponse
func UserToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		School:    user.School,
		StudentID: user.StudentID,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ProblemToResponse 将Problem模型转换为ProblemResponse
func ProblemToResponse(problem *models.Problem) ProblemResponse {
	var tags []string
	if problem.Tags != "" {
		json.Unmarshal([]byte(problem.Tags), &tags)
	}

	response := ProblemResponse{
		ID:          problem.ID,
		Title:       problem.Title,
		Description: problem.Description,
		InputDesc:   problem.InputDesc,
		OutputDesc:  problem.OutputDesc,
		Difficulty:  problem.Difficulty,
		Category:    problem.Category,
		Tags:        tags,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		IsPublic:    problem.IsPublic,
		AuthorID:    problem.AuthorID,
		SubmitCount: problem.SubmitCount,
		AcceptCount: problem.AcceptedCount,
		CreatedAt:   problem.CreatedAt,
		UpdatedAt:   problem.UpdatedAt,
	}

	// 注意：测试用例需要单独查询和转换

	return response
}

// TestCaseToResponse 将TestCase模型转换为TestCaseResponse
func TestCaseToResponse(testCase *models.TestCase) TestCaseResponse {
	return TestCaseResponse{
		ID:        testCase.ID,
		ProblemID: testCase.ProblemID,
		Input:     testCase.Input,
		Output:    testCase.Output,
		IsSample:  testCase.IsSample,
		CreatedAt: testCase.CreatedAt,
		UpdatedAt: testCase.UpdatedAt,
	}
}

// SubmissionToResponse 将Submission模型转换为SubmissionResponse
func SubmissionToResponse(submission *models.Submission) SubmissionResponse {
	return SubmissionResponse{
		ID:           submission.ID,
		UserID:       submission.UserID,
		ProblemID:    submission.ProblemID,
		Code:         submission.Code,
		Language:     submission.Language,
		Status:       submission.Status,
		Score:        submission.Score,
		RunTime:      submission.RunTime,
		Memory:       submission.Memory,
		ErrorMessage: submission.ErrorMessage,
		PassedTests:  submission.PassedTests,
		TotalTests:   submission.TotalTests,
		JudgeID:      submission.JudgeID,
		CreatedAt:    submission.CreatedAt,
		UpdatedAt:    submission.UpdatedAt,
	}
}

// ForumPostToResponse 将ForumPost模型转换为ForumPostResponse
func ForumPostToResponse(post *models.ForumPost) ForumPostResponse {
	var tags []string
	if post.Tags != "" {
		json.Unmarshal([]byte(post.Tags), &tags)
	}

	response := ForumPostResponse{
		ID:         post.ID,
		Title:      post.Title,
		Content:    post.Content,
		Category:   post.Category,
		Tags:       tags,
		ViewCount:  post.ViewCount,
		ReplyCount: post.ReplyCount,
		LikeCount:  post.LikeCount,
		IsLocked:   post.IsLocked,
		IsSticky:   post.IsSticky,
		IsPublic:   post.IsPublic,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}

	// 处理用户信息
	if post.User.ID != 0 {
		response.User = UserResponse{
			ID:       post.User.ID,
			Username: post.User.Username,
			Email:    post.User.Email,
			Nickname: post.User.Nickname,
			Avatar:   post.User.Avatar,
			Role:     post.User.Role,
		}
	}

	return response
}

// ForumReplyToResponse 将ForumReply模型转换为ForumReplyResponse
func ForumReplyToResponse(reply *models.ForumReply) ForumReplyResponse {
	return ForumReplyResponse{
		ID:        reply.ID,
		Content:   reply.Content,
		PostID:    reply.PostID,
		UserID:    reply.UserID,
		ParentID:  reply.ParentID,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}
}

// NewsToResponse 将News模型转换为NewsResponse
func NewsToResponse(news *models.News) NewsResponse {
	var tags []string
	if news.Tags != "" {
		json.Unmarshal([]byte(news.Tags), &tags)
	}

	return NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		Summary:     news.Summary,
		Category:    news.Category,
		Tags:        tags,
		IsPublished: news.IsPublished,
		AuthorID:    news.AuthorID,
		ViewCount:   news.ViewCount,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	}
}

// ProblemsToResponse 批量转换Problem模型为ProblemResponse
func ProblemsToResponse(problems []models.Problem) []ProblemResponse {
	responses := make([]ProblemResponse, len(problems))
	for i, problem := range problems {
		responses[i] = ProblemToResponse(&problem)
	}
	return responses
}

// SubmissionsToResponse 批量转换Submission模型为SubmissionResponse
func SubmissionsToResponse(submissions []models.Submission) []SubmissionResponse {
	responses := make([]SubmissionResponse, len(submissions))
	for i, submission := range submissions {
		responses[i] = SubmissionToResponse(&submission)
	}
	return responses
}

// ForumPostsToResponse 批量转换ForumPost模型为ForumPostResponse
func ForumPostsToResponse(posts []models.ForumPost) []ForumPostResponse {
	responses := make([]ForumPostResponse, len(posts))
	for i, post := range posts {
		responses[i] = ForumPostToResponse(&post)
	}
	return responses
}

// ForumRepliesToResponse 批量转换ForumReply模型为ForumReplyResponse
func ForumRepliesToResponse(replies []models.ForumReply) []ForumReplyResponse {
	responses := make([]ForumReplyResponse, len(replies))
	for i, reply := range replies {
		responses[i] = ForumReplyToResponse(&reply)
	}
	return responses
}

// NewsListToResponse 批量转换News模型为NewsResponse
func NewsListToResponse(newsList []models.News) []NewsResponse {
	responses := make([]NewsResponse, len(newsList))
	for i, news := range newsList {
		responses[i] = NewsToResponse(&news)
	}
	return responses
}

// ========== DTO ↔ Domain 转换函数 ==========

// UserRegisterRequestToDomain 将UserRegisterRequest转换为Domain实体
func UserRegisterRequestToDomain(req *UserRegisterRequest) *domain.User {
	return &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Nickname:  req.Nickname,
		School:    req.School,
		StudentID: req.StudentID,
		Role:      "student", // 默认角色
	}
}

// UserDomainToResponse 将Domain实体转换为UserResponse
func UserDomainToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		School:    user.School,
		StudentID: user.StudentID,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ProblemCreateRequestToDomain 将ProblemCreateRequest转换为Domain实体
func ProblemCreateRequestToDomain(req *ProblemCreateRequest) *domain.Problem {
	return &domain.Problem{
		Title:       req.Title,
		Description: req.Description,
		InputDesc:   req.InputDesc,
		OutputDesc:  req.OutputDesc,
		Difficulty:  req.Difficulty,
		Category:    req.Category,
		Tags:        req.Tags,
		TimeLimit:   req.TimeLimit,
		MemoryLimit: req.MemoryLimit,
		IsPublic:    false, // 默认私有
	}
}

// ProblemDomainToResponse 将Domain实体转换为ProblemResponse
func ProblemDomainToResponse(problem *domain.Problem) ProblemResponse {
	return ProblemResponse{
		ID:          problem.ID,
		Title:       problem.Title,
		Description: problem.Description,
		InputDesc:   problem.InputDesc,
		OutputDesc:  problem.OutputDesc,
		Difficulty:  problem.Difficulty,
		Category:    problem.Category,
		Tags:        problem.Tags,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		IsPublic:    problem.IsPublic,
		AuthorID:    problem.AuthorID,
		SubmitCount: problem.SubmitCount,
		AcceptCount: problem.AcceptedCount,
		CreatedAt:   problem.CreatedAt,
		UpdatedAt:   problem.UpdatedAt,
	}
}

// SubmissionCreateRequestToDomain 将SubmissionCreateRequest转换为Domain实体
func SubmissionCreateRequestToDomain(req *SubmissionCreateRequest) *domain.Submission {
	return &domain.Submission{
		ProblemID: req.ProblemID,
		Code:      req.Code,
		Language:  req.Language,
		Status:    "pending", // 默认状态
	}
}

// SubmissionDomainToResponse 将Domain实体转换为SubmissionResponse
func SubmissionDomainToResponse(submission *domain.Submission) SubmissionResponse {
	return SubmissionResponse{
		ID:           submission.ID,
		UserID:       submission.UserID,
		ProblemID:    submission.ProblemID,
		Code:         submission.Code,
		Language:     submission.Language,
		Status:       submission.Status,
		Score:        submission.Score,
		ErrorMessage: submission.ErrorMessage,
		PassedTests:  submission.PassedTests,
		TotalTests:   submission.TotalTests,
		CreatedAt:    submission.CreatedAt,
		UpdatedAt:    submission.UpdatedAt,
	}
}

// ForumPostCreateRequestToDomain 将ForumPostCreateRequest转换为Domain实体
func ForumPostCreateRequestToDomain(req *ForumPostCreateRequest) *domain.ForumPost {
	return &domain.ForumPost{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
		IsPublic: true, // 默认公开
	}
}

// ForumPostDomainToResponse 将Domain实体转换为ForumPostResponse
func ForumPostDomainToResponse(post *domain.ForumPost) ForumPostResponse {
	response := ForumPostResponse{
		ID:         post.ID,
		Title:      post.Title,
		Content:    post.Content,
		Category:   post.Category,
		Tags:       post.Tags,
		ViewCount:  post.ViewCount,
		ReplyCount: post.ReplyCount,
		LikeCount:  post.LikeCount,
		IsLocked:   post.IsLocked,
		IsSticky:   post.IsSticky,
		IsPublic:   post.IsPublic,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}

	// 处理用户信息
	if post.User != nil {
		response.User = UserResponse{
			ID:       post.User.ID,
			Username: post.User.Username,
			Email:    post.User.Email,
			Nickname: post.User.Nickname,
			Avatar:   post.User.Avatar,
			Role:     post.User.Role,
		}
	}

	return response
}

// ForumReplyCreateRequestToDomain 将ForumReplyCreateRequest转换为Domain实体
func ForumReplyCreateRequestToDomain(req *ForumReplyCreateRequest, postID uint) *domain.ForumReply {
	return &domain.ForumReply{
		PostID:   postID,
		Content:  req.Content,
		ParentID: req.ParentID,
	}
}

// ForumReplyDomainToResponse 将Domain实体转换为ForumReplyResponse
func ForumReplyDomainToResponse(reply *domain.ForumReply) ForumReplyResponse {
	return ForumReplyResponse{
		ID:        reply.ID,
		Content:   reply.Content,
		PostID:    reply.PostID,
		UserID:    reply.AuthorID,
		ParentID:  reply.ParentID,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}
}

// NewsCreateRequestToDomain 将NewsCreateRequest转换为Domain实体
func NewsCreateRequestToDomain(req *NewsCreateRequest) *domain.News {
	status := "draft"
	if req.IsPublished {
		status = "published"
	}
	return &domain.News{
		Title:    req.Title,
		Content:  req.Content,
		Summary:  req.Summary,
		Category: req.Category,
		Tags:     req.Tags,
		Status:   status,
	}
}

// NewsDomainToResponse 将Domain实体转换为NewsResponse
func NewsDomainToResponse(news *domain.News) NewsResponse {
	return NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		Summary:     news.Summary,
		Category:    news.Category,
		Tags:        news.Tags,
		IsPublished: news.Status == "published",
		AuthorID:    news.AuthorID,
		ViewCount:   news.ViewCount,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	}
}

// TestCaseDomainToResponse 将Domain实体转换为TestCaseResponse
func TestCaseDomainToResponse(testCase *domain.TestCase) TestCaseResponse {
	return TestCaseResponse{
		ID:        testCase.ID,
		ProblemID: testCase.ProblemID,
		Input:     testCase.Input,
		Output:    testCase.Output,
		IsSample:  testCase.IsSample,
		CreatedAt: testCase.CreatedAt,
		UpdatedAt: testCase.UpdatedAt,
	}
}
