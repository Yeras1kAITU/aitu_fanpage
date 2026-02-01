package service

import (
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/repository"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
	stats    map[primitive.ObjectID]int
	statsMu  sync.RWMutex
	likeChan chan likeRequest
}

type likeRequest struct {
	postID primitive.ObjectID
	userID primitive.ObjectID
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) *PostService {
	ps := &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
		stats:    make(map[primitive.ObjectID]int),
		likeChan: make(chan likeRequest, 100),
	}

	return ps
}

func (s *PostService) CreatePost(req dto.CreatePostRequest, authorID primitive.ObjectID) (*models.Post, error) {
	user, err := s.userRepo.FindByID(authorID)
	var authorName string

	if err != nil {
		authorName = "Anonymous User"
	} else {
		authorName = user.DisplayName
	}

	post := models.NewPost(
		req.Title,
		req.Content,
		req.Description,
		models.PostCategory(req.Category),
		authorID,
		authorName,
	)

	for i, media := range req.Media {
		if i >= 10 {
			break
		}
		post.AddMedia(media.URL, media.Type, media.Caption)
	}

	if err := s.postRepo.Create(post); err != nil {
		return nil, err
	}

	go s.updatePostStats(post.ID)

	return post, nil
}

func (s *PostService) updatePostStats(postID primitive.ObjectID) {
	time.Sleep(2 * time.Second)

	s.statsMu.Lock()
	s.stats[postID] = 0
	s.statsMu.Unlock()
}
