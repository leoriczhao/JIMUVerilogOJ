package services

import (
	"errors"
	"strconv"
	"time"
	"verilog-oj/backend/internal/domain"
)

// SubmissionRepository 提交仓储接口
type SubmissionRepository interface {
	// 创建提交
	Create(submission *domain.Submission) error
	// 根据ID获取提交
	GetByID(id uint) (*domain.Submission, error)
	// 获取提交列表
	List(page, limit int, userID, problemID uint, status string) ([]domain.Submission, int64, error)
	// 更新提交状态
	UpdateStatus(id uint, status string, score int, runTime, memory int, errorMessage string, passedTests, totalTests int) error
	// 统计用户通过的题目数
	CountAcceptedByUser(userID, problemID uint) (int64, error)
	// 获取提交统计信息
	GetStats(userID uint) (map[string]interface{}, error)
	// 软删除提交记录
	SoftDelete(id uint) error
}

// SubmissionService 提交服务
type SubmissionService struct {
	submissionRepo SubmissionRepository
	problemRepo    ProblemRepository
	userRepo       UserRepository
}

// NewSubmissionService 创建提交服务
func NewSubmissionService(submissionRepo SubmissionRepository, problemRepo ProblemRepository, userRepo UserRepository) *SubmissionService {
	return &SubmissionService{
		submissionRepo: submissionRepo,
		problemRepo:    problemRepo,
		userRepo:       userRepo,
	}
}

// SubmissionListResult 提交列表结果（内部使用）
type SubmissionListResult struct {
	Submissions []domain.Submission
	Total       int64
	Page        int
	Limit       int
}

// CreateSubmission 创建提交
func (s *SubmissionService) CreateSubmission(problemID uint, code, language string, userID uint) (*domain.Submission, error) {
	// 验证用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 验证题目是否存在
	problem, err := s.problemRepo.GetByID(problemID)
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return nil, errors.New("题目不存在")
	}

	// 验证题目是否公开
	if !problem.IsPublic {
		// 检查用户是否有权限访问私有题目
		if problem.AuthorID != userID {
			return nil, errors.New("没有权限访问此题目")
		}
	}

	// 设置默认语言
	if language == "" {
		language = "verilog"
	}

	// 验证代码长度
	if len(code) > 100000 { // 100KB限制
		return nil, errors.New("代码长度超过限制")
	}

	// 验证代码不能为空
	if code == "" {
		return nil, errors.New("代码不能为空")
	}

	// 创建提交记录
	submission := &domain.Submission{
		UserID:    userID,
		ProblemID: problemID,
		Code:      code,
		Language:  language,
		Status:    "pending",
		Score:     0,
	}

	if err := s.submissionRepo.Create(submission); err != nil {
		return nil, err
	}

	// 更新题目的提交统计
	if err := s.problemRepo.UpdateSubmitCount(problemID, 1); err != nil {
		// 记录错误但不影响提交创建
	}

	// 更新用户的提交统计
	if err := s.userRepo.UpdateStats(userID, user.Solved, user.Submitted+1); err != nil {
		// 记录错误但不影响提交创建
	}

	return submission, nil
}

// GetSubmission 获取提交详情
func (s *SubmissionService) GetSubmission(id uint) (*domain.Submission, error) {
	submission, err := s.submissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if submission == nil {
		return nil, errors.New("提交记录不存在")
	}

	return submission, nil
}

// ListSubmissions 获取提交列表
func (s *SubmissionService) ListSubmissions(page, limit int, userID, problemID uint, status string) (*SubmissionListResult, error) {
	// 验证分页参数
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	submissions, total, err := s.submissionRepo.List(page, limit, userID, problemID, status)
	if err != nil {
		return nil, err
	}

	return &SubmissionListResult{
		Submissions: submissions,
		Total:       total,
		Page:        page,
		Limit:       limit,
	}, nil
}

// UpdateSubmissionStatus 更新提交状态
func (s *SubmissionService) UpdateSubmissionStatus(id uint, status string, score int, runTime, memory int, errorMessage string, passedTests, totalTests int) error {
	// 更新提交状态
	if err := s.submissionRepo.UpdateStatus(id, status, score, runTime, memory, errorMessage, passedTests, totalTests); err != nil {
		return err
	}

	// 如果状态是accepted，更新用户和题目的统计
	if status == "accepted" {
		// 获取提交信息
		submission, err := s.submissionRepo.GetByID(id)
		if err != nil {
			return err
		}
		if submission == nil {
			return errors.New("提交记录不存在")
		}

		// 检查是否是第一次通过此题
		existingAccepted, err := s.submissionRepo.CountAcceptedByUser(submission.UserID, submission.ProblemID)
		if err != nil {
			return err
		}

		if existingAccepted == 1 { // 这是第一次通过
			// 获取用户当前信息
			user, err := s.userRepo.GetByID(submission.UserID)
			if err != nil {
				return err
			}
			if user != nil {
				// 更新用户的解题数
				if err := s.userRepo.UpdateStats(submission.UserID, user.Solved+1, user.Submitted); err != nil {
					return err
				}
			}

			// 更新题目的通过数
			if err := s.problemRepo.UpdateAcceptedCount(submission.ProblemID, 1); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUserSubmissions 获取用户的提交记录
func (s *SubmissionService) GetUserSubmissions(userID uint, page, limit int) (*SubmissionListResult, error) {
	return s.ListSubmissions(page, limit, userID, 0, "")
}

// GetProblemSubmissions 获取题目的提交记录
func (s *SubmissionService) GetProblemSubmissions(problemID uint, page, limit int) (*SubmissionListResult, error) {
	return s.ListSubmissions(page, limit, 0, problemID, "")
}

// GetSubmissionStats 获取提交统计信息
func (s *SubmissionService) GetSubmissionStats(userID uint) (map[string]interface{}, error) {
	return s.submissionRepo.GetStats(userID)
}

// ValidateSubmissionAccess 验证用户是否有权限访问提交记录
func (s *SubmissionService) ValidateSubmissionAccess(submissionID, userID uint) error {
	submission, err := s.submissionRepo.GetByID(submissionID)
	if err != nil {
		return err
	}
	if submission == nil {
		return errors.New("提交记录不存在")
	}

	// 检查是否是提交者本人或管理员
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	if submission.UserID != userID && user.Role != "admin" && user.Role != "teacher" {
		return errors.New("没有权限访问此提交记录")
	}

	return nil
}

// DeleteSubmission 删除提交记录（仅管理员或提交者本人）
func (s *SubmissionService) DeleteSubmission(id uint, userID uint, userRole string) error {
	submission, err := s.GetSubmission(id)
	if err != nil {
		return err
	}

	// 权限检查：只有提交者本人或管理员可以删除
	if submission.UserID != userID && userRole != "admin" {
		return errors.New("没有权限删除此提交记录")
	}

	// 执行软删除
	if err := s.submissionRepo.SoftDelete(id); err != nil {
		return err
	}

	return nil
}

// generateJudgeID 生成判题ID
func generateJudgeID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
