package service

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/repository"
)

type PostService struct {
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	commentRepo repository.CommentRepository
	rateLimiter *RateLimiter
	likeTracker *LikeTracker
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository, commentRepo repository.CommentRepository) *PostService {
	service := &PostService{
		postRepo:    postRepo,
		userRepo:    userRepo,
		commentRepo: commentRepo,
		rateLimiter: NewRateLimiter(),
		likeTracker: NewLikeTracker(),
	}

	// Start background cleanup
	go service.backgroundCleanup()

	return service
}

func (s *PostService) CreatePost(req dto.CreatePostRequest, authorID primitive.ObjectID, uploadedFiles []*UploadedFile) (*models.Post, error) {
	user, err := s.userRepo.FindByID(authorID)
	if err != nil {
		return nil, err
	}

	if !user.CanCreatePost() {
		return nil, errors.New("user cannot create posts")
	}

	authorName := user.DisplayName
	if authorName == "" {
		authorName = "Anonymous User"
	}

	// Create post
	post := models.NewPost(
		req.Title,
		req.Content,
		req.Description,
		models.PostCategory(req.Category),
		authorID,
		authorName,
	)

	if len(req.Tags) > 0 {
		post.AddTags(req.Tags...)
	}

	for i, uploadedFile := range uploadedFiles {
		if i >= 10 {
			break
		}

		caption := ""
		if i < len(req.Media) {
			caption = req.Media[i].Caption
		}

		post.AddMedia(
			uploadedFile.URL,
			string(uploadedFile.MediaType),
			caption,
			uploadedFile.FileSize,
			uploadedFile.ThumbnailURL,
			uploadedFile.Dimensions,
			uploadedFile.Checksum,
		)
	}

	post.CalculatePopularityScore()

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	user.IncrementPostCount()
	s.userRepo.Update(user)

	return post, nil
}

func (s *PostService) GetPosts(limit, offset int) ([]*models.Post, error) {
	return s.postRepo.FindAll(limit, offset)
}

func (s *PostService) GetFeed(userID primitive.ObjectID, category string, limit, offset int) ([]*models.Post, error) {
	if category != "" {
		return s.postRepo.FindByCategory(category, limit)
	}
	return s.postRepo.FindAll(limit, offset)
}

func (s *PostService) GetPostsByCategory(category string, limit int) ([]*models.Post, error) {
	return s.postRepo.FindByCategory(category, limit)
}

func (s *PostService) GetPostsByAuthor(authorID primitive.ObjectID, limit int) ([]*models.Post, error) {
	return s.postRepo.FindByAuthor(authorID, limit)
}

func (s *PostService) GetPostByID(postID primitive.ObjectID) (*models.Post, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}

	go s.postRepo.IncrementViewCount(postID)
	post.IncrementViewCount()

	return post, nil
}

func (s *PostService) GetPinnedPosts(limit int) ([]*models.Post, error) {
	return s.postRepo.FindPinned(limit)
}

func (s *PostService) GetFeaturedPosts(limit int) ([]*models.Post, error) {
	return s.postRepo.FindFeatured(limit)
}

func (s *PostService) GetPopularPosts(limit int, days int) ([]*models.Post, error) {
	return s.postRepo.FindPopular(limit, days)
}

func (s *PostService) GetPostsByTags(tags []string, limit int) ([]*models.Post, error) {
	return s.postRepo.FindByTags(tags, limit)
}

func (s *PostService) SearchPosts(query string, limit int) ([]*models.Post, error) {
	return s.postRepo.Search(query, limit)
}

func (s *PostService) GetCategoriesStats() (map[string]int, error) {
	return s.postRepo.GetCategoriesStats()
}

func (s *PostService) LikePost(postID, userID primitive.ObjectID) error {
	if !s.rateLimiter.CanLike(userID) {
		return errors.New("rate limit exceeded: maximum 3 likes per 3 minutes")
	}

	s.rateLimiter.RecordLike(userID)

	timestamp := time.Now()
	s.likeTracker.AddLike(postID, userID, timestamp)

	if err := s.postRepo.IncrementLikeCount(postID); err != nil {
		return errors.New("failed to update like count: " + err.Error())
	}

	return nil
}

func (s *PostService) UnlikePost(postID, userID primitive.ObjectID) error {
	likeRecord := s.likeTracker.CanUnlike(postID, userID)
	if likeRecord == nil {
		return errors.New("cannot unlike: no active likes on this post within the last 3 minutes")
	}

	if !s.likeTracker.Unlike(likeRecord.ID) {
		return errors.New("unlike failed")
	}

	if err := s.postRepo.DecrementLikeCount(postID); err != nil {
		return errors.New("failed to update database: " + err.Error())
	}

	return nil
}

func (s *PostService) GetUserLikeStats(userID primitive.ObjectID) (int, time.Time, error) {
	return 0, time.Time{}, nil
}

func (s *PostService) backgroundCleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.rateLimiter.Cleanup()
		s.likeTracker.Cleanup()
	}
}

func (s *PostService) UpdatePost(postID, userID primitive.ObjectID, req dto.UpdatePostRequest) (*models.Post, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.CanEditPost(post.AuthorID) {
		return nil, errors.New("not authorized to edit this post")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Description != "" {
		post.Description = req.Description
	}
	if req.Category != "" {
		post.Category = models.PostCategory(req.Category)
	}
	if len(req.Tags) > 0 {
		post.Tags = req.Tags
	}

	post.CalculatePopularityScore()

	if err := s.postRepo.Update(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(postID, userID primitive.ObjectID) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !user.CanDeletePost(post.AuthorID) {
		return errors.New("not authorized to delete this post")
	}

	if err := s.commentRepo.DeleteByPostID(postID); err != nil {
		return err
	}

	return s.postRepo.Delete(postID)
}

func (s *PostService) PinPost(postID, userID primitive.ObjectID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !user.CanManagePosts() {
		return errors.New("not authorized to pin posts")
	}

	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}

	post.Pin()
	return s.postRepo.Update(post)
}

func (s *PostService) UnpinPost(postID, userID primitive.ObjectID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !user.CanManagePosts() {
		return errors.New("not authorized to unpin posts")
	}

	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}

	post.Unpin()
	return s.postRepo.Update(post)
}

func (s *PostService) FeaturePost(postID, userID primitive.ObjectID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !user.CanManagePosts() {
		return errors.New("not authorized to feature posts")
	}

	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}

	post.Feature()
	return s.postRepo.Update(post)
}

func (s *PostService) UnfeaturePost(postID, userID primitive.ObjectID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !user.CanManagePosts() {
		return errors.New("not authorized to unfeature posts")
	}

	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}

	post.Unfeature()
	return s.postRepo.Update(post)
}

func (s *PostService) GetPostLikeCount(postID primitive.ObjectID) (int, error) {
	count := s.likeTracker.GetPostLikeCount(postID)
	return count, nil
}

func (s *PostService) HasUserLiked(postID, userID primitive.ObjectID) bool {
	return false
}

func (s *PostService) GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}
