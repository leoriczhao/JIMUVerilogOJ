package services

import (
	"errors"
	"log"
	"verilog-oj/backend/internal/domain"
)

// ForumRepository 论坛仓储接口
type ForumRepository interface {
	// 帖子相关
	CreatePost(post *domain.ForumPost) error
	GetPostByID(id uint) (*domain.ForumPost, error)
	ListPosts(page, limit int, filters map[string]interface{}) ([]domain.ForumPost, int64, error)
	UpdatePost(post *domain.ForumPost) error
	DeletePost(id uint) error
	IncrementPostViewCount(id uint) error
	IncrementPostReplyCount(id uint) error
	DecrementPostReplyCount(id uint) error

	// 回复相关
	CreateReply(reply *domain.ForumReply) error
	GetRepliesByPostID(postID uint, page, limit int) ([]domain.ForumReply, int64, error)
	UpdateReply(reply *domain.ForumReply) error
	DeleteReply(id uint) error
	DeleteRepliesByPostID(postID uint) error

	// 点赞相关
	CreateLike(like *domain.ForumLike) error
	DeleteLike(userID uint, postID, replyID *uint) error
	CheckLikeExists(userID uint, postID, replyID *uint) (bool, error)
	GetLikeCount(postID, replyID *uint) (int64, error)
	GetUserLikes(userID uint, page, limit int) ([]domain.ForumLike, int64, error)
}

// ForumService 论坛服务
type ForumService struct {
	forumRepo ForumRepository
	userRepo  UserRepository
}

// NewForumService 创建论坛服务
func NewForumService(forumRepo ForumRepository, userRepo UserRepository) *ForumService {
	return &ForumService{
		forumRepo: forumRepo,
		userRepo:  userRepo,
	}
}

// CreatePost 创建帖子
func (s *ForumService) CreatePost(post *domain.ForumPost) error {
	// 业务逻辑验证
	if post.Title == "" {
		return errors.New("帖子标题不能为空")
	}
	if post.Content == "" {
		return errors.New("帖子内容不能为空")
	}
	if post.AuthorID == 0 {
		return errors.New("用户ID不能为空")
	}

	// 验证用户是否存在
	user, err := s.userRepo.GetByID(post.AuthorID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	return s.forumRepo.CreatePost(post)
}

// GetPost 获取帖子详情
func (s *ForumService) GetPost(id uint) (*domain.ForumPost, error) {
	post, err := s.forumRepo.GetPostByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("帖子不存在")
	}

	// 更新浏览量
	if err := s.forumRepo.IncrementPostViewCount(id); err != nil {
		log.Printf("failed to increment post view count: %v", err)
	}

	return post, nil
}

// ListPosts 获取帖子列表
func (s *ForumService) ListPosts(page, limit int, filters map[string]interface{}) ([]domain.ForumPost, int64, error) {
	// 验证分页参数
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// 默认只显示公开帖子
	if filters == nil {
		filters = make(map[string]interface{})
	}
	if _, exists := filters["is_public"]; !exists {
		filters["is_public"] = true
	}

	return s.forumRepo.ListPosts(page, limit, filters)
}

// UpdatePost 更新帖子
func (s *ForumService) UpdatePost(post *domain.ForumPost) error {
	// 业务逻辑验证
	if post.Title == "" {
		return errors.New("帖子标题不能为空")
	}
	if post.Content == "" {
		return errors.New("帖子内容不能为空")
	}

	return s.forumRepo.UpdatePost(post)
}

// DeletePost 删除帖子
func (s *ForumService) DeletePost(id uint) error {
	// 检查帖子是否存在
	post, err := s.forumRepo.GetPostByID(id)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("帖子不存在")
	}

	// 删除帖子的所有回复
	if err := s.forumRepo.DeleteRepliesByPostID(id); err != nil {
		return err
	}

	// 删除帖子
	return s.forumRepo.DeletePost(id)
}

// CreateReply 创建回复
func (s *ForumService) CreateReply(reply *domain.ForumReply) error {
	// 业务逻辑验证
	if reply.Content == "" {
		return errors.New("回复内容不能为空")
	}
	if reply.AuthorID == 0 {
		return errors.New("用户ID不能为空")
	}
	if reply.PostID == 0 {
		return errors.New("帖子ID不能为空")
	}

	// 验证用户是否存在
	user, err := s.userRepo.GetByID(reply.AuthorID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 验证帖子是否存在
	post, err := s.forumRepo.GetPostByID(reply.PostID)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("帖子不存在")
	}

	// 创建回复
	if err := s.forumRepo.CreateReply(reply); err != nil {
		return err
	}

	// 增加帖子回复数
	if err := s.forumRepo.IncrementPostReplyCount(reply.PostID); err != nil {
		log.Printf("failed to increment post reply count: %v", err)
	}

	return nil
}

// ListReplies 获取回复列表
func (s *ForumService) ListReplies(postID uint, page, limit int) ([]domain.ForumReply, int64, error) {
	// 验证分页参数
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// 验证帖子是否存在
	post, err := s.forumRepo.GetPostByID(postID)
	if err != nil {
		return nil, 0, err
	}
	if post == nil {
		return nil, 0, errors.New("帖子不存在")
	}

	return s.forumRepo.GetRepliesByPostID(postID, page, limit)
}

// ToggleLike 切换点赞状态
func (s *ForumService) ToggleLike(userID uint, postID, replyID *uint) error {
	// 验证用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 检查是否已点赞
	exists, err := s.forumRepo.CheckLikeExists(userID, postID, replyID)
	if err != nil {
		return err
	}

	if exists {
		// 取消点赞
		return s.forumRepo.DeleteLike(userID, postID, replyID)
	} else {
		// 添加点赞
		like := &domain.ForumLike{
			UserID: userID,
		}
		if postID != nil {
			like.TargetType = "post"
			like.TargetID = *postID
		} else if replyID != nil {
			like.TargetType = "reply"
			like.TargetID = *replyID
		}
		return s.forumRepo.CreateLike(like)
	}
}

// GetLikeCount 获取点赞数
func (s *ForumService) GetLikeCount(postID, replyID *uint) (int64, error) {
	return s.forumRepo.GetLikeCount(postID, replyID)
}

// CheckUserLike 检查用户是否已点赞
func (s *ForumService) CheckUserLike(userID uint, postID, replyID *uint) (bool, error) {
	return s.forumRepo.CheckLikeExists(userID, postID, replyID)
}
