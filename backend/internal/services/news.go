package services

import (
	"errors"
	"verilog-oj/backend/internal/domain"
)

// NewsRepository 新闻仓储接口
type NewsRepository interface {
	// 创建新闻
	Create(news *domain.News) error
	// 根据ID获取新闻
	GetByID(id uint) (*domain.News, error)
	// 获取新闻列表
	List(page, limit int, filters map[string]interface{}) ([]domain.News, int64, error)
	// 更新新闻
	Update(news *domain.News) error
	// 删除新闻
	Delete(id uint) error
	// 增加浏览量
	IncrementViewCount(id uint) error
}

// NewsService 新闻服务
type NewsService struct {
	newsRepo NewsRepository
	userRepo UserRepository
}

// NewNewsService 创建新闻服务
func NewNewsService(newsRepo NewsRepository, userRepo UserRepository) *NewsService {
	return &NewsService{
		newsRepo: newsRepo,
		userRepo: userRepo,
	}
}

// CreateNews 创建新闻
func (s *NewsService) CreateNews(news *domain.News) error {
	// 业务逻辑验证
	if news.Title == "" {
		return errors.New("新闻标题不能为空")
	}
	if news.Content == "" {
		return errors.New("新闻内容不能为空")
	}
	if news.AuthorID == 0 {
		return errors.New("作者ID不能为空")
	}

	// 验证作者是否存在
	user, err := s.userRepo.GetByID(news.AuthorID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("作者不存在")
	}

	return s.newsRepo.Create(news)
}

// GetNewsList 获取新闻列表
func (s *NewsService) GetNewsList(page, limit int) ([]domain.News, int64, error) {
	// 验证分页参数
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// 默认只显示已发布的新闻
	filters := make(map[string]interface{})
	filters["is_published"] = true

	return s.newsRepo.List(page, limit, filters)
}

// GetNews 获取新闻详情
func (s *NewsService) GetNews(id uint) (*domain.News, error) {
	news, err := s.newsRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if news == nil {
		return nil, errors.New("新闻不存在")
	}

	// 更新浏览量
	if err := s.newsRepo.IncrementViewCount(id); err != nil {
		// 记录错误但不影响获取新闻
	}

	return news, nil
}

// GetNewsByID 根据ID获取新闻详情（handlers接口兼容方法）
func (s *NewsService) GetNewsByID(id uint) (*domain.News, error) {
	return s.GetNews(id)
}

// UpdateNews 更新新闻
func (s *NewsService) UpdateNews(news *domain.News) error {
	// 业务逻辑验证
	if news.Title == "" {
		return errors.New("新闻标题不能为空")
	}
	if news.Content == "" {
		return errors.New("新闻内容不能为空")
	}

	return s.newsRepo.Update(news)
}

// DeleteNews 删除新闻
func (s *NewsService) DeleteNews(id uint) error {
	// 检查新闻是否存在
	news, err := s.newsRepo.GetByID(id)
	if err != nil {
		return err
	}
	if news == nil {
		return errors.New("新闻不存在")
	}

	return s.newsRepo.Delete(id)
}
