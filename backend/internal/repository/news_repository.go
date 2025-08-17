package repository

import (
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/services"

	"gorm.io/gorm"
)

// NewsRepository 新闻仓储实现
type NewsRepository struct {
	db *gorm.DB
}

// NewNewsRepository 创建新闻仓储实例
func NewNewsRepository(db *gorm.DB) services.NewsRepository {
	return &NewsRepository{
		db: db,
	}
}

// Create 创建新闻
func (r *NewsRepository) Create(news *domain.News) error {
	modelNews := NewsDomainToModel(news)
	err := r.db.Create(modelNews).Error
	if err != nil {
		return err
	}

	// 更新ID和时间戳
	news.ID = modelNews.ID
	news.CreatedAt = modelNews.CreatedAt
	news.UpdatedAt = modelNews.UpdatedAt

	return nil
}

// GetByID 根据ID获取新闻
func (r *NewsRepository) GetByID(id uint) (*domain.News, error) {
	var modelNews models.News
	err := r.db.First(&modelNews, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return NewsModelToDomain(&modelNews), nil
}

// List 获取新闻列表
func (r *NewsRepository) List(page, limit int, filters map[string]interface{}) ([]domain.News, int64, error) {
	var modelNewsList []models.News
	var total int64

	query := r.db.Model(&models.News{})

	// 应用过滤条件
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("is_featured DESC, created_at DESC").Find(&modelNewsList).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为domain对象
	newsList := make([]domain.News, len(modelNewsList))
	for i, modelNews := range modelNewsList {
		newsList[i] = *NewsModelToDomain(&modelNews)
	}

	return newsList, total, nil
}

// Update 更新新闻
func (r *NewsRepository) Update(news *domain.News) error {
	modelNews := NewsDomainToModel(news)
	err := r.db.Save(modelNews).Error
	if err != nil {
		return err
	}

	// 更新时间戳
	news.UpdatedAt = modelNews.UpdatedAt

	return nil
}

// Delete 删除新闻
func (r *NewsRepository) Delete(id uint) error {
	return r.db.Delete(&models.News{}, id).Error
}

// IncrementViewCount 增加浏览量
func (r *NewsRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.News{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
