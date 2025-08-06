package repository

import (
	"encoding/json"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/services"

	"gorm.io/gorm"
)

// ForumRepository 论坛仓储实现
type ForumRepository struct {
	db *gorm.DB
}

// NewForumRepository 创建论坛仓储实例
func NewForumRepository(db *gorm.DB) services.ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

// 转换函数：domain -> model
func (r *ForumRepository) domainPostToModel(post *domain.ForumPost) *models.ForumPost {
	var tagsJSON string
	if len(post.Tags) > 0 {
		tagsBytes, _ := json.Marshal(post.Tags)
		tagsJSON = string(tagsBytes)
	}

	modelPost := &models.ForumPost{
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
	}
	
	// 只有在更新时才设置ID和时间戳
	if post.ID != 0 {
		modelPost.ID = post.ID
		modelPost.CreatedAt = post.CreatedAt
		modelPost.UpdatedAt = post.UpdatedAt
	}
	
	return modelPost
}

// 转换函数：model -> domain
func (r *ForumRepository) modelPostToDomain(post *models.ForumPost) *domain.ForumPost {
	var tags []string
	if post.Tags != "" {
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

func (r *ForumRepository) domainReplyToModel(reply *domain.ForumReply) *models.ForumReply {
	modelReply := &models.ForumReply{
		Content:   reply.Content,
		PostID:    reply.PostID,
		UserID:    reply.AuthorID,
		ParentID:  reply.ParentID,
		LikeCount: reply.LikeCount,
	}
	
	// 只有在更新时才设置ID和时间戳
	if reply.ID != 0 {
		modelReply.ID = reply.ID
		modelReply.CreatedAt = reply.CreatedAt
		modelReply.UpdatedAt = reply.UpdatedAt
	}
	
	return modelReply
}

func (r *ForumRepository) modelReplyToDomain(reply *models.ForumReply) *domain.ForumReply {
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

func (r *ForumRepository) domainLikeToModel(like *domain.ForumLike) *models.ForumLike {
	modelLike := &models.ForumLike{
		ID:        like.ID,
		UserID:    like.UserID,
		CreatedAt: like.CreatedAt,
	}

	switch like.TargetType {
	case "post":
		modelLike.PostID = &like.TargetID
	case "reply":
		modelLike.ReplyID = &like.TargetID
	}

	return modelLike
}

func (r *ForumRepository) modelLikeToDomain(like *models.ForumLike) *domain.ForumLike {
	domainLike := &domain.ForumLike{
		ID:        like.ID,
		UserID:    like.UserID,
		CreatedAt: like.CreatedAt,
	}

	if like.PostID != nil {
		domainLike.TargetType = "post"
		domainLike.TargetID = *like.PostID
	} else if like.ReplyID != nil {
		domainLike.TargetType = "reply"
		domainLike.TargetID = *like.ReplyID
	}

	return domainLike
}

// 帖子相关方法
func (r *ForumRepository) CreatePost(post *domain.ForumPost) error {
	modelPost := r.domainPostToModel(post)
	if err := r.db.Create(modelPost).Error; err != nil {
		return err
	}
	// 将生成的ID回写到domain对象中
	post.ID = modelPost.ID
	post.CreatedAt = modelPost.CreatedAt
	post.UpdatedAt = modelPost.UpdatedAt
	return nil
}

func (r *ForumRepository) GetPostByID(id uint) (*domain.ForumPost, error) {
	var modelPost models.ForumPost
	err := r.db.Preload("User").First(&modelPost, id).Error
	if err != nil {
		return nil, err
	}
	return r.modelPostToDomain(&modelPost), nil
}

func (r *ForumRepository) ListPosts(page, limit int, filters map[string]interface{}) ([]domain.ForumPost, int64, error) {
	var modelPosts []models.ForumPost
	var total int64

	query := r.db.Model(&models.ForumPost{}).Preload("User")

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
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&modelPosts).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为domain对象
	domainPosts := make([]domain.ForumPost, len(modelPosts))
	for i, modelPost := range modelPosts {
		domainPosts[i] = *r.modelPostToDomain(&modelPost)
	}

	return domainPosts, total, nil
}

func (r *ForumRepository) UpdatePost(post *domain.ForumPost) error {
	modelPost := r.domainPostToModel(post)
	return r.db.Save(modelPost).Error
}

func (r *ForumRepository) DeletePost(id uint) error {
	return r.db.Delete(&models.ForumPost{}, id).Error
}

func (r *ForumRepository) IncrementPostViewCount(id uint) error {
	return r.db.Model(&models.ForumPost{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *ForumRepository) IncrementPostReplyCount(id uint) error {
	return r.db.Model(&models.ForumPost{}).Where("id = ?", id).Update("reply_count", gorm.Expr("reply_count + 1")).Error
}

func (r *ForumRepository) DecrementPostReplyCount(id uint) error {
	return r.db.Model(&models.ForumPost{}).Where("id = ?", id).Update("reply_count", gorm.Expr("reply_count - 1")).Error
}

// 回复相关方法
func (r *ForumRepository) CreateReply(reply *domain.ForumReply) error {
	modelReply := r.domainReplyToModel(reply)
	if err := r.db.Create(modelReply).Error; err != nil {
		return err
	}
	// 将生成的ID回写到domain对象中
	reply.ID = modelReply.ID
	reply.CreatedAt = modelReply.CreatedAt
	reply.UpdatedAt = modelReply.UpdatedAt
	return nil
}

func (r *ForumRepository) GetRepliesByPostID(postID uint, page, limit int) ([]domain.ForumReply, int64, error) {
	var modelReplies []models.ForumReply
	var total int64

	query := r.db.Model(&models.ForumReply{}).Where("post_id = ?", postID).Preload("User")

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at ASC").Find(&modelReplies).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为domain对象
	domainReplies := make([]domain.ForumReply, len(modelReplies))
	for i, modelReply := range modelReplies {
		domainReplies[i] = *r.modelReplyToDomain(&modelReply)
	}

	return domainReplies, total, nil
}

func (r *ForumRepository) UpdateReply(reply *domain.ForumReply) error {
	modelReply := r.domainReplyToModel(reply)
	return r.db.Save(modelReply).Error
}

func (r *ForumRepository) DeleteReply(id uint) error {
	return r.db.Delete(&models.ForumReply{}, id).Error
}

func (r *ForumRepository) DeleteRepliesByPostID(postID uint) error {
	return r.db.Where("post_id = ?", postID).Delete(&models.ForumReply{}).Error
}

// 点赞相关方法
func (r *ForumRepository) CreateLike(like *domain.ForumLike) error {
	modelLike := r.domainLikeToModel(like)
	return r.db.Create(modelLike).Error
}

func (r *ForumRepository) DeleteLike(userID uint, postID, replyID *uint) error {
	query := r.db.Where("user_id = ?", userID)

	if postID != nil {
		query = query.Where("post_id = ? AND reply_id IS NULL", *postID)
	} else if replyID != nil {
		query = query.Where("reply_id = ? AND post_id IS NULL", *replyID)
	} else {
		return gorm.ErrInvalidData
	}

	return query.Delete(&models.ForumLike{}).Error
}

func (r *ForumRepository) CheckLikeExists(userID uint, postID, replyID *uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.ForumLike{}).Where("user_id = ?", userID)

	if postID != nil {
		query = query.Where("post_id = ? AND reply_id IS NULL", *postID)
	} else if replyID != nil {
		query = query.Where("reply_id = ? AND post_id IS NULL", *replyID)
	} else {
		return false, gorm.ErrInvalidData
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *ForumRepository) GetLikeCount(postID, replyID *uint) (int64, error) {
	var count int64
	query := r.db.Model(&models.ForumLike{})

	if postID != nil {
		query = query.Where("post_id = ? AND reply_id IS NULL", *postID)
	} else if replyID != nil {
		query = query.Where("reply_id = ? AND post_id IS NULL", *replyID)
	} else {
		return 0, gorm.ErrInvalidData
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *ForumRepository) GetUserLikes(userID uint, page, limit int) ([]domain.ForumLike, int64, error) {
	var modelLikes []models.ForumLike
	var total int64

	query := r.db.Model(&models.ForumLike{}).Where("user_id = ?", userID).Preload("Post").Preload("Reply")

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&modelLikes).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为domain对象
	domainLikes := make([]domain.ForumLike, len(modelLikes))
	for i, modelLike := range modelLikes {
		domainLikes[i] = *r.modelLikeToDomain(&modelLike)
	}

	return domainLikes, total, nil
}
