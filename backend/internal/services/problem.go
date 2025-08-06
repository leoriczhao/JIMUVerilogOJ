package services

import (
	"errors"
	"verilog-oj/backend/internal/domain"
)

// ProblemRepository 题目仓储接口
type ProblemRepository interface {
	// 创建题目
	Create(problem *domain.Problem) error
	// 根据ID获取题目
	GetByID(id uint) (*domain.Problem, error)
	// 获取题目列表
	List(page, limit int, filters map[string]interface{}) ([]domain.Problem, int64, error)
	// 更新题目
	Update(problem *domain.Problem) error
	// 删除题目
	Delete(id uint) error
	// 更新题目提交统计
	UpdateSubmitCount(id uint, increment int) error
	// 更新题目通过统计
	UpdateAcceptedCount(id uint, increment int) error
	// 获取题目的测试用例
	GetTestCases(problemID uint) ([]domain.TestCase, error)
	// 创建测试用例
	CreateTestCase(testCase *domain.TestCase) error
	// 删除题目的所有测试用例
	DeleteTestCases(problemID uint) error
}

// ProblemService 题目服务
type ProblemService struct {
	problemRepo ProblemRepository
}

// NewProblemService 创建题目服务
func NewProblemService(problemRepo ProblemRepository) *ProblemService {
	return &ProblemService{
		problemRepo: problemRepo,
	}
}

// CreateProblem 创建题目
func (s *ProblemService) CreateProblem(problem *domain.Problem) error {
	// 业务逻辑验证
	if problem.Title == "" {
		return errors.New("题目标题不能为空")
	}
	if problem.Description == "" {
		return errors.New("题目描述不能为空")
	}
	if problem.TimeLimit <= 0 {
		return errors.New("时间限制必须大于0")
	}
	if problem.MemoryLimit <= 0 {
		return errors.New("内存限制必须大于0")
	}

	return s.problemRepo.Create(problem)
}

// GetProblem 获取题目详情
func (s *ProblemService) GetProblem(id uint) (*domain.Problem, error) {
	problem, err := s.problemRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return nil, errors.New("题目不存在")
	}
	return problem, nil
}

// ListProblems 获取题目列表
func (s *ProblemService) ListProblems(page, limit int, filters map[string]interface{}) ([]domain.Problem, int64, error) {
	// 验证分页参数
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return s.problemRepo.List(page, limit, filters)
}

// UpdateProblem 更新题目
func (s *ProblemService) UpdateProblem(problem *domain.Problem) error {
	// 业务逻辑验证
	if problem.Title == "" {
		return errors.New("题目标题不能为空")
	}
	if problem.Description == "" {
		return errors.New("题目描述不能为空")
	}

	return s.problemRepo.Update(problem)
}

// DeleteProblem 删除题目
func (s *ProblemService) DeleteProblem(id uint) error {
	// 检查题目是否存在
	problem, err := s.problemRepo.GetByID(id)
	if err != nil {
		return err
	}
	if problem == nil {
		return errors.New("题目不存在")
	}

	// 删除相关的测试用例
	if err := s.problemRepo.DeleteTestCases(id); err != nil {
		return err
	}

	// 删除题目
	return s.problemRepo.Delete(id)
}

// GetTestCases 获取题目的测试用例
func (s *ProblemService) GetTestCases(problemID uint) ([]domain.TestCase, error) {
	// 检查题目是否存在
	problem, err := s.problemRepo.GetByID(problemID)
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return nil, errors.New("题目不存在")
	}

	return s.problemRepo.GetTestCases(problemID)
}

// AddTestCase 添加测试用例
func (s *ProblemService) AddTestCase(testCase *domain.TestCase) error {
	// 验证测试用例
	if testCase.ProblemID == 0 {
		return errors.New("题目ID不能为空")
	}
	if testCase.Input == "" && testCase.Output == "" {
		return errors.New("测试用例输入和输出不能同时为空")
	}

	// 检查题目是否存在
	problem, err := s.problemRepo.GetByID(testCase.ProblemID)
	if err != nil {
		return err
	}
	if problem == nil {
		return errors.New("题目不存在")
	}

	return s.problemRepo.CreateTestCase(testCase)
}

// UpdateProblemStats 更新题目统计信息
func (s *ProblemService) UpdateProblemStats(problemID uint, submitIncrement, acceptedIncrement int) error {
	if submitIncrement != 0 {
		if err := s.problemRepo.UpdateSubmitCount(problemID, submitIncrement); err != nil {
			return err
		}
	}
	if acceptedIncrement != 0 {
		if err := s.problemRepo.UpdateAcceptedCount(problemID, acceptedIncrement); err != nil {
			return err
		}
	}
	return nil
}
