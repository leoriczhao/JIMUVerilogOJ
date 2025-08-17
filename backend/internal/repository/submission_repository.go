package repository

import (
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/services"

	"gorm.io/gorm"
)

// SubmissionRepository 提交仓储实现
type SubmissionRepository struct {
	db *gorm.DB
}

// NewSubmissionRepository 创建提交仓储实例
func NewSubmissionRepository(db *gorm.DB) services.SubmissionRepository {
	return &SubmissionRepository{
		db: db,
	}
}

// Create 创建提交
func (r *SubmissionRepository) Create(submission *domain.Submission) error {
	modelSubmission := SubmissionDomainToModel(submission)
	err := r.db.Create(modelSubmission).Error
	if err != nil {
		return err
	}

	// 更新ID和时间戳
	submission.ID = modelSubmission.ID
	submission.CreatedAt = modelSubmission.CreatedAt
	submission.UpdatedAt = modelSubmission.UpdatedAt

	return nil
}

// GetByID 根据ID获取提交
func (r *SubmissionRepository) GetByID(id uint) (*domain.Submission, error) {
	var modelSubmission models.Submission
	err := r.db.First(&modelSubmission, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return SubmissionModelToDomain(&modelSubmission), nil
}

// List 获取提交列表
func (r *SubmissionRepository) List(page, limit int, userID, problemID uint, status string) ([]domain.Submission, int64, error) {
	var modelSubmissions []models.Submission
	var total int64

	query := r.db.Model(&models.Submission{})

	// 应用过滤条件
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if problemID > 0 {
		query = query.Where("problem_id = ?", problemID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&modelSubmissions).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为domain对象
	submissions := make([]domain.Submission, len(modelSubmissions))
	for i, modelSubmission := range modelSubmissions {
		submissions[i] = *SubmissionModelToDomain(&modelSubmission)
	}

	return submissions, total, nil
}

// UpdateStatus 更新提交状态
func (r *SubmissionRepository) UpdateStatus(id uint, status string, score int, runTime, memory int, errorMessage string, passedTests, totalTests int) error {
	updates := map[string]interface{}{
		"status":        status,
		"score":         score,
		"run_time":      runTime,
		"memory":        memory,
		"error_message": errorMessage,
		"passed_tests":  passedTests,
		"total_tests":   totalTests,
	}

	return r.db.Model(&models.Submission{}).Where("id = ?", id).Updates(updates).Error
}

// CountAcceptedByUser 统计用户通过的题目数
func (r *SubmissionRepository) CountAcceptedByUser(userID, problemID uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.Submission{}).Where("user_id = ? AND status = ?", userID, "Accepted")

	if problemID > 0 {
		query = query.Where("problem_id = ?", problemID)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetStats 获取提交统计信息
func (r *SubmissionRepository) GetStats(userID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总提交数
	var totalSubmissions int64
	err := r.db.Model(&models.Submission{}).Where("user_id = ?", userID).Count(&totalSubmissions).Error
	if err != nil {
		return nil, err
	}
	stats["total_submissions"] = totalSubmissions

	// 通过的提交数
	var acceptedSubmissions int64
	err = r.db.Model(&models.Submission{}).Where("user_id = ? AND status = ?", userID, "Accepted").Count(&acceptedSubmissions).Error
	if err != nil {
		return nil, err
	}
	stats["accepted_submissions"] = acceptedSubmissions

	// 通过的题目数（去重）
	var solvedProblems int64
	err = r.db.Model(&models.Submission{}).Where("user_id = ? AND status = ?", userID, "Accepted").Distinct("problem_id").Count(&solvedProblems).Error
	if err != nil {
		return nil, err
	}
	stats["solved_problems"] = solvedProblems

	// 计算通过率
	if totalSubmissions > 0 {
		stats["acceptance_rate"] = float64(acceptedSubmissions) / float64(totalSubmissions) * 100
	} else {
		stats["acceptance_rate"] = 0.0
	}

	return stats, nil
}

// SoftDelete 软删除提交记录
func (r *SubmissionRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.Submission{}, id).Error
}
