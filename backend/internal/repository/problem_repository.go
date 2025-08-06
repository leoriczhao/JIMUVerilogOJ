package repository

import (
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/services"

	"gorm.io/gorm"
)

// ProblemRepository 题目仓储实现
type ProblemRepository struct {
	db *gorm.DB
}

// NewProblemRepository 创建题目仓储实例
func NewProblemRepository(db *gorm.DB) services.ProblemRepository {
	return &ProblemRepository{
		db: db,
	}
}

// Create 创建题目
func (r *ProblemRepository) Create(problem *domain.Problem) error {
	modelProblem := ProblemDomainToModel(problem)
	err := r.db.Create(modelProblem).Error
	if err != nil {
		return err
	}

	// 更新ID和时间戳
	problem.ID = modelProblem.ID
	problem.CreatedAt = modelProblem.CreatedAt
	problem.UpdatedAt = modelProblem.UpdatedAt

	return nil
}

// GetByID 根据ID获取题目
func (r *ProblemRepository) GetByID(id uint) (*domain.Problem, error) {
	var modelProblem models.Problem
	err := r.db.First(&modelProblem, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ProblemModelToDomain(&modelProblem), nil
}

// List 获取题目列表
func (r *ProblemRepository) List(page, limit int, filters map[string]interface{}) ([]domain.Problem, int64, error) {
	var modelProblems []models.Problem
	var total int64

	query := r.db.Model(&models.Problem{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "difficulty":
			if difficulty, ok := value.(string); ok && difficulty != "" {
				query = query.Where("difficulty = ?", difficulty)
			}
		case "category":
			if category, ok := value.(string); ok && category != "" {
				query = query.Where("category = ?", category)
			}
		case "title":
			if title, ok := value.(string); ok && title != "" {
				query = query.Where("title LIKE ?", "%"+title+"%")
			}
		case "is_public":
			query = query.Where("is_public = ?", value)
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&modelProblems).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为domain对象
	problems := make([]domain.Problem, len(modelProblems))
	for i, modelProblem := range modelProblems {
		problems[i] = *ProblemModelToDomain(&modelProblem)
	}

	return problems, total, nil
}

// Update 更新题目
func (r *ProblemRepository) Update(problem *domain.Problem) error {
	modelProblem := ProblemDomainToModel(problem)
	err := r.db.Save(modelProblem).Error
	if err != nil {
		return err
	}

	// 更新时间戳
	problem.UpdatedAt = modelProblem.UpdatedAt

	return nil
}

// Delete 删除题目
func (r *ProblemRepository) Delete(id uint) error {
	return r.db.Delete(&models.Problem{}, id).Error
}

// CreateTestCase 创建测试用例
func (r *ProblemRepository) CreateTestCase(testCase *domain.TestCase) error {
	modelTestCase := TestCaseDomainToModel(testCase)
	err := r.db.Create(modelTestCase).Error
	if err != nil {
		return err
	}

	// 更新ID和时间戳
	testCase.ID = modelTestCase.ID
	testCase.CreatedAt = modelTestCase.CreatedAt
	testCase.UpdatedAt = modelTestCase.UpdatedAt

	return nil
}

// GetTestCases 获取题目的测试用例
func (r *ProblemRepository) GetTestCases(problemID uint) ([]domain.TestCase, error) {
	var modelTestCases []models.TestCase
	err := r.db.Where("problem_id = ?", problemID).Order("id ASC").Find(&modelTestCases).Error
	if err != nil {
		return nil, err
	}

	// 转换为domain对象
	testCases := make([]domain.TestCase, len(modelTestCases))
	for i, modelTestCase := range modelTestCases {
		testCases[i] = *TestCaseModelToDomain(&modelTestCase)
	}

	return testCases, nil
}

// DeleteTestCases 删除题目的测试用例
func (r *ProblemRepository) DeleteTestCases(problemID uint) error {
	return r.db.Where("problem_id = ?", problemID).Delete(&models.TestCase{}).Error
}

// UpdateSubmitCount 更新题目提交统计
func (r *ProblemRepository) UpdateSubmitCount(id uint, increment int) error {
	return r.db.Model(&models.Problem{}).Where("id = ?", id).Update("submit_count", gorm.Expr("submit_count + ?", increment)).Error
}

// UpdateAcceptedCount 更新题目通过统计
func (r *ProblemRepository) UpdateAcceptedCount(id uint, increment int) error {
	return r.db.Model(&models.Problem{}).Where("id = ?", id).Update("accepted_count", gorm.Expr("accepted_count + ?", increment)).Error
}
