package repository_test

import (
	"testing"
	"verilog-oj/backend/internal/domain"
	"verilog-oj/backend/internal/models"
	"verilog-oj/backend/internal/repository"
	"verilog-oj/backend/internal/services"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupForumTestDB creates a new in-memory SQLite database for testing the forum repository.
func setupForumTestDB(t *testing.T) (*gorm.DB, services.UserRepository) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to memory db: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.ForumPost{}, &models.ForumReply{}, &models.ForumLike{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	return db, userRepo
}

func TestForumRepository_CreateAndGetPost(t *testing.T) {
    db, userRepo := setupForumTestDB(t)
    repo := repository.NewForumRepository(db)

    // 1. Create a user first, as posts require an author.
    user := &domain.User{Username: "post_author", Email: "author@test.com", Password: "password"}
    err := userRepo.Create(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)

    // 2. Create a post
    post := &domain.ForumPost{
        Title:    "Test Post Title",
        Content:  "This is the content of the test post.",
        AuthorID: user.ID,
        Category: "general",
        Tags:     []string{"go", "testing"},
    }
    err = repo.CreatePost(post)
    assert.NoError(t, err)
    assert.NotZero(t, post.ID)

    // 3. Verify the post can be retrieved from the DB
    retrievedPost, err := repo.GetPostByID(post.ID)
    assert.NoError(t, err)
    assert.NotNil(t, retrievedPost)
    assert.Equal(t, "Test Post Title", retrievedPost.Title)
    assert.Equal(t, user.ID, retrievedPost.AuthorID)
    assert.Equal(t, []string{"go", "testing"}, retrievedPost.Tags)
    // Check that the associated user is preloaded correctly
    assert.NotNil(t, retrievedPost.User)
    assert.Equal(t, user.Username, retrievedPost.User.Username)
}

func TestForumRepository_UpdatePost(t *testing.T) {
	db, userRepo := setupForumTestDB(t)
	repo := repository.NewForumRepository(db)

	// 1. Create a user and an initial post
	user := &domain.User{Username: "updater", Email: "updater@test.com", Password: "password"}
	err := userRepo.Create(user)
	assert.NoError(t, err)
	post := &domain.ForumPost{
		Title: "Original Title", Content: "Original Content", AuthorID: user.ID,
	}
	err = repo.CreatePost(post)
	assert.NoError(t, err)

	// 2. Update the post
	post.Title = "Updated Title"
	post.Content = "Updated Content"
	post.Tags = []string{"updated"}
	err = repo.UpdatePost(post)
	assert.NoError(t, err)

	// 3. Verify the update
	retrievedPost, err := repo.GetPostByID(post.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", retrievedPost.Title)
	assert.Equal(t, "Updated Content", retrievedPost.Content)
	assert.Equal(t, []string{"updated"}, retrievedPost.Tags)
}

func TestForumRepository_ListPosts(t *testing.T) {
	db, userRepo := setupForumTestDB(t)
	repo := repository.NewForumRepository(db)

	// Create users and posts
	user1 := &domain.User{Username: "user1", Email: "user1@test.com", Password: "pw"}
	userRepo.Create(user1)
	user2 := &domain.User{Username: "user2", Email: "user2@test.com", Password: "pw"}
	userRepo.Create(user2)

	repo.CreatePost(&domain.ForumPost{Title: "Post 1", AuthorID: user1.ID, Category: "cat1"})
	repo.CreatePost(&domain.ForumPost{Title: "Post 2", AuthorID: user2.ID, Category: "cat2"})
	repo.CreatePost(&domain.ForumPost{Title: "Post 3", AuthorID: user1.ID, Category: "cat1"})

	// Test case 1: List all posts
	posts, total, err := repo.ListPosts(1, 10, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, posts, 3)

	// Test case 2: Filter by category
	posts, total, err = repo.ListPosts(1, 10, map[string]interface{}{"category": "cat1"})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, posts, 2)

	// Test case 3: Test pagination
	posts, total, err = repo.ListPosts(2, 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, posts, 1)
}

func TestForumRepository_DeletePost(t *testing.T) {
	db, userRepo := setupForumTestDB(t)
	repo := repository.NewForumRepository(db)

	user := &domain.User{Username: "deleter", Email: "deleter@test.com", Password: "pw"}
	userRepo.Create(user)
	post := &domain.ForumPost{Title: "To Be Deleted", AuthorID: user.ID}
	repo.CreatePost(post)
	assert.NotZero(t, post.ID)

	err := repo.DeletePost(post.ID)
	assert.NoError(t, err)

	_, err = repo.GetPostByID(post.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestForumRepository_PostCounters(t *testing.T) {
	db, userRepo := setupForumTestDB(t)
	repo := repository.NewForumRepository(db)

	user := &domain.User{Username: "counter", Email: "counter@test.com", Password: "pw"}
	userRepo.Create(user)
	post := &domain.ForumPost{Title: "Counter Post", AuthorID: user.ID}
	repo.CreatePost(post)

	// Test View Count
	err := repo.IncrementPostViewCount(post.ID)
	assert.NoError(t, err)
	retrieved, _ := repo.GetPostByID(post.ID)
	assert.Equal(t, 1, retrieved.ViewCount)

	// Test Reply Count
	err = repo.IncrementPostReplyCount(post.ID)
	assert.NoError(t, err)
	err = repo.IncrementPostReplyCount(post.ID)
	assert.NoError(t, err)
	retrieved, _ = repo.GetPostByID(post.ID)
	assert.Equal(t, 2, retrieved.ReplyCount)

	err = repo.DecrementPostReplyCount(post.ID)
	assert.NoError(t, err)
	retrieved, _ = repo.GetPostByID(post.ID)
	assert.Equal(t, 1, retrieved.ReplyCount)
}

func TestForumRepository_CreateAndGetReplies(t *testing.T) {
	db, userRepo := setupForumTestDB(t)
	repo := repository.NewForumRepository(db)

	// Setup user and post
	user := &domain.User{Username: "replier", Email: "replier@test.com", Password: "pw"}
	userRepo.Create(user)
	post := &domain.ForumPost{Title: "Post with Replies", AuthorID: user.ID}
	repo.CreatePost(post)

	// Create replies
	reply1 := &domain.ForumReply{Content: "Reply 1", PostID: post.ID, AuthorID: user.ID}
	repo.CreateReply(reply1)
	assert.NotZero(t, reply1.ID)

	reply2 := &domain.ForumReply{Content: "Reply 2", PostID: post.ID, AuthorID: user.ID, ParentID: &reply1.ID}
	repo.CreateReply(reply2)
	assert.NotZero(t, reply2.ID)


	// Get replies
	replies, total, err := repo.GetRepliesByPostID(post.ID, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, replies, 2)
	assert.Equal(t, "Reply 1", replies[0].Content)
	assert.Equal(t, "Reply 2", replies[1].Content)
	assert.NotNil(t, replies[1].ParentID)
	assert.Equal(t, reply1.ID, *replies[1].ParentID)
}

func TestForumRepository_Likes(t *testing.T) {
	db, userRepo := setupForumTestDB(t)
	repo := repository.NewForumRepository(db)

	// Setup user and post
	user1 := &domain.User{Username: "liker1", Email: "liker1@test.com", Password: "pw"}
	userRepo.Create(user1)
	user2 := &domain.User{Username: "liker2", Email: "liker2@test.com", Password: "pw"}
	userRepo.Create(user2)
	post := &domain.ForumPost{Title: "Likable Post", AuthorID: user1.ID}
	repo.CreatePost(post)

	// Test Like Creation
	like := &domain.ForumLike{UserID: user1.ID, TargetID: post.ID, TargetType: "post"}
	err := repo.CreateLike(like)
	assert.NoError(t, err)

	// Test CheckLikeExists
	exists, err := repo.CheckLikeExists(user1.ID, &post.ID, nil)
	assert.NoError(t, err)
	assert.True(t, exists)
	exists, err = repo.CheckLikeExists(user2.ID, &post.ID, nil)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test GetLikeCount
	count, err := repo.GetLikeCount(&post.ID, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Test Like Deletion
	err = repo.DeleteLike(user1.ID, &post.ID, nil)
	assert.NoError(t, err)
	exists, err = repo.CheckLikeExists(user1.ID, &post.ID, nil)
	assert.NoError(t, err)
	assert.False(t, exists)
	count, err = repo.GetLikeCount(&post.ID, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
