package repository

import (
	"encoding/json"
	"time"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"

	"gorm.io/gorm"
)

// ========== Domain ↔ Model 转换函数 ==========

// UserDomainToModel 将Domain实体转换为Model
func UserDomainToModel(user *domain.User) *models.User {
	return &models.User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		School:      user.School,
		StudentID:   user.StudentID,
		Role:        user.Role,
		Solved:      user.Solved,
		Submitted:   user.Submitted,
		Rating:      user.Rating,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// UserModelToDomain 将Model转换为Domain实体
func UserModelToDomain(user *models.User) *domain.User {
	return &domain.User{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		School:      user.School,
		StudentID:   user.StudentID,
		Role:        user.Role,
		Solved:      user.Solved,
		Submitted:   user.Submitted,
		Rating:      user.Rating,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// ProblemDomainToModel 将Domain实体转换为Model
func ProblemDomainToModel(problem *domain.Problem) *models.Problem {
	// 将Tags切片转换为JSON字符串
	tagsJSON := "[]"
	if len(problem.Tags) > 0 {
		if tagsBytes, err := json.Marshal(problem.Tags); err == nil {
			tagsJSON = string(tagsBytes)
		}
	}

	return &models.Problem{
		ID:            problem.ID,
		Title:         problem.Title,
		Description:   problem.Description,
		InputDesc:     problem.InputDesc,
		OutputDesc:    problem.OutputDesc,
		Difficulty:    problem.Difficulty,
		Category:      problem.Category,
		Tags:          tagsJSON,
		TimeLimit:     problem.TimeLimit,
		MemoryLimit:   problem.MemoryLimit,
		SubmitCount:   problem.SubmitCount,
		AcceptedCount: problem.AcceptedCount,
		IsPublic:      problem.IsPublic,
		AuthorID:      problem.AuthorID,
		CreatedAt:     problem.CreatedAt,
		UpdatedAt:     problem.UpdatedAt,
	}
}

// ProblemModelToDomain 将Model转换为Domain实体
func ProblemModelToDomain(problem *models.Problem) *domain.Problem {
	// 将JSON字符串转换为Tags切片
	var tags []string
	if problem.Tags != "" && problem.Tags != "[]" {
		json.Unmarshal([]byte(problem.Tags), &tags)
	}

	return &domain.Problem{
		ID:            problem.ID,
		Title:         problem.Title,
		Description:   problem.Description,
		InputDesc:     problem.InputDesc,
		OutputDesc:    problem.OutputDesc,
		Difficulty:    problem.Difficulty,
		Category:      problem.Category,
		Tags:          tags,
		TimeLimit:     problem.TimeLimit,
		MemoryLimit:   problem.MemoryLimit,
		SubmitCount:   problem.SubmitCount,
		AcceptedCount: problem.AcceptedCount,
		IsPublic:      problem.IsPublic,
		AuthorID:      problem.AuthorID,
		CreatedAt:     problem.CreatedAt,
		UpdatedAt:     problem.UpdatedAt,
	}
}

// SubmissionDomainToModel 将Domain实体转换为Model
func SubmissionDomainToModel(submission *domain.Submission) *models.Submission {
	return &models.Submission{
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

// SubmissionModelToDomain 将Model转换为Domain实体
func SubmissionModelToDomain(submission *models.Submission) *domain.Submission {
	return &domain.Submission{
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

// ForumPostDomainToModel 将Domain实体转换为Model
func ForumPostDomainToModel(post *domain.ForumPost) *models.ForumPost {
	// 将Tags切片转换为JSON字符串
	tagsJSON := "[]"
	if len(post.Tags) > 0 {
		if tagsBytes, err := json.Marshal(post.Tags); err == nil {
			tagsJSON = string(tagsBytes)
		}
	}

	return &models.ForumPost{
		ID:         post.ID,
		Title:      post.Title,
		Content:    post.Content,
		UserID:     post.AuthorID,
		Category:   post.Category,
		Tags:       tagsJSON,
		ViewCount:  post.ViewCount,
		ReplyCount: post.ReplyCount,
		LikeCount:  post.LikeCount,
		IsSticky:   post.IsSticky,
		IsLocked:   post.IsLocked,
		IsPublic:   post.IsPublic,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}

// ForumPostModelToDomain 将Model转换为Domain实体
func ForumPostModelToDomain(post *models.ForumPost) *domain.ForumPost {
	// 将JSON字符串转换为Tags切片
	var tags []string
	if post.Tags != "" && post.Tags != "[]" {
		json.Unmarshal([]byte(post.Tags), &tags)
	}

	return &domain.ForumPost{
		ID:         post.ID,
		Title:      post.Title,
		Content:    post.Content,
		AuthorID:   post.UserID,
		Category:   post.Category,
		Tags:       tags,
		ViewCount:  post.ViewCount,
		ReplyCount: post.ReplyCount,
		LikeCount:  post.LikeCount,
		IsSticky:   post.IsSticky,
		IsLocked:   post.IsLocked,
		IsPublic:   post.IsPublic,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}

// ForumReplyDomainToModel 将Domain实体转换为Model
func ForumReplyDomainToModel(reply *domain.ForumReply) *models.ForumReply {
	return &models.ForumReply{
		ID:        reply.ID,
		Content:   reply.Content,
		PostID:    reply.PostID,
		UserID:    reply.AuthorID,
		ParentID:  reply.ParentID,
		LikeCount: reply.LikeCount,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}
}

// ForumReplyModelToDomain 将Model转换为Domain实体
func ForumReplyModelToDomain(reply *models.ForumReply) *domain.ForumReply {
	return &domain.ForumReply{
		ID:        reply.ID,
		Content:   reply.Content,
		PostID:    reply.PostID,
		AuthorID:  reply.UserID,
		ParentID:  reply.ParentID,
		LikeCount: reply.LikeCount,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}
}

// NewsDomainToModel 将Domain实体转换为Model
func NewsDomainToModel(news *domain.News) *models.News {
	// 将Tags切片转换为JSON字符串
	tagsJSON := "[]"
	if len(news.Tags) > 0 {
		if tagsBytes, err := json.Marshal(news.Tags); err == nil {
			tagsJSON = string(tagsBytes)
		}
	}

	// 根据状态设置IsPublished
	isPublished := news.Status == "published"

	return &models.News{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		Summary:     news.Summary,
		AuthorID:    news.AuthorID,
		IsPublished: isPublished,
		IsFeatured:  news.IsFeatured,
		Category:    news.Category,
		Tags:        tagsJSON,
		ViewCount:   news.ViewCount,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	}
}

// NewsModelToDomain 将Model转换为Domain实体
func NewsModelToDomain(news *models.News) *domain.News {
	// 将JSON字符串转换为Tags切片
	var tags []string
	if news.Tags != "" && news.Tags != "[]" {
		json.Unmarshal([]byte(news.Tags), &tags)
	}

	// 根据IsPublished设置状态
	status := "draft"
	if news.IsPublished {
		status = "published"
	}

	// 设置PublishedAt
	var publishedAt *time.Time
	if news.IsPublished {
		publishedAt = &news.UpdatedAt
	}

	return &domain.News{
		ID:          news.ID,
		Title:       news.Title,
		Content:     news.Content,
		Summary:     news.Summary,
		AuthorID:    news.AuthorID,
		Status:      status,
		Category:    news.Category,
		Tags:        tags,
		IsFeatured:  news.IsFeatured,
		ViewCount:   news.ViewCount,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
		PublishedAt: publishedAt,
	}
}

// TestCaseDomainToModel 将Domain实体转换为Model
func TestCaseDomainToModel(testCase *domain.TestCase) *models.TestCase {
	return &models.TestCase{
		ID:        testCase.ID,
		ProblemID: testCase.ProblemID,
		Input:     testCase.Input,
		Output:    testCase.Output,
		IsSample:  testCase.IsSample,
		CreatedAt: testCase.CreatedAt,
		UpdatedAt: testCase.UpdatedAt,
	}
}

// TestCaseModelToDomain 将Model转换为Domain实体
func TestCaseModelToDomain(testCase *models.TestCase) *domain.TestCase {
	return &domain.TestCase{
		ID:        testCase.ID,
		ProblemID: testCase.ProblemID,
		Input:     testCase.Input,
		Output:    testCase.Output,
		IsSample:  testCase.IsSample,
		CreatedAt: testCase.CreatedAt,
		UpdatedAt: testCase.UpdatedAt,
	}
}

// ========== 批量转换函数 ==========

// UsersModelToDomain 批量转换User Model为Domain
func UsersModelToDomain(users []models.User) []domain.User {
	result := make([]domain.User, len(users))
	for i, user := range users {
		result[i] = *UserModelToDomain(&user)
	}
	return result
}

// ProblemsModelToDomain 批量转换Problem Model为Domain
func ProblemsModelToDomain(problems []models.Problem) []domain.Problem {
	result := make([]domain.Problem, len(problems))
	for i, problem := range problems {
		result[i] = *ProblemModelToDomain(&problem)
	}
	return result
}

// SubmissionsModelToDomain 批量转换Submission Model为Domain
func SubmissionsModelToDomain(submissions []models.Submission) []domain.Submission {
	result := make([]domain.Submission, len(submissions))
	for i, submission := range submissions {
		result[i] = *SubmissionModelToDomain(&submission)
	}
	return result
}

// ForumPostsModelToDomain 批量转换ForumPost Model为Domain
func ForumPostsModelToDomain(posts []models.ForumPost) []domain.ForumPost {
	result := make([]domain.ForumPost, len(posts))
	for i, post := range posts {
		result[i] = *ForumPostModelToDomain(&post)
	}
	return result
}

// ForumRepliesModelToDomain 批量转换ForumReply Model为Domain
func ForumRepliesModelToDomain(replies []models.ForumReply) []domain.ForumReply {
	result := make([]domain.ForumReply, len(replies))
	for i, reply := range replies {
		result[i] = *ForumReplyModelToDomain(&reply)
	}
	return result
}

// NewsListModelToDomain 批量转换News Model为Domain
func NewsListModelToDomain(newsList []models.News) []domain.News {
	result := make([]domain.News, len(newsList))
	for i, news := range newsList {
		result[i] = *NewsModelToDomain(&news)
	}
	return result
}

// TestCasesModelToDomain 批量转换TestCase Model为Domain
func TestCasesModelToDomain(testCases []models.TestCase) []domain.TestCase {
	result := make([]domain.TestCase, len(testCases))
	for i, testCase := range testCases {
		result[i] = *TestCaseModelToDomain(&testCase)
	}
	return result
}

// ========== 辅助函数 ==========

// SetModelTimestamps 设置Model的时间戳字段
func SetModelTimestamps(model interface{}) {
	now := time.Now()
	switch m := model.(type) {
	case *models.User:
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
	case *models.Problem:
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
	case *models.Submission:
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
	case *models.ForumPost:
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
	case *models.ForumReply:
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
	case *models.News:
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
	case *models.TestCase:
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		m.UpdatedAt = now
	}
}

// HandleSoftDelete 处理软删除
func HandleSoftDelete(db *gorm.DB, model interface{}) error {
	now := time.Now()
	switch m := model.(type) {
	case *models.User:
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	case *models.Problem:
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	case *models.Submission:
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	case *models.ForumPost:
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	case *models.ForumReply:
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	case *models.News:
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	case *models.TestCase:
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}
	return db.Save(model).Error
}
