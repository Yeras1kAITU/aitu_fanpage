package service

import (
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LikeRecord tracks a single like
type LikeRecord struct {
	ID        primitive.ObjectID
	UserID    primitive.ObjectID
	PostID    primitive.ObjectID
	Timestamp time.Time
	Active    bool
}

// LikeTracker tracks all likes with ability to unlike
type LikeTracker struct {
	mu    sync.RWMutex
	likes map[primitive.ObjectID][]LikeRecord
}

func NewLikeTracker() *LikeTracker {
	return &LikeTracker{
		likes: make(map[primitive.ObjectID][]LikeRecord),
	}
}

// AddLike adds new like record
func (lt *LikeTracker) AddLike(postID, userID primitive.ObjectID, timestamp time.Time) primitive.ObjectID {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	likeID := primitive.NewObjectID()
	record := LikeRecord{
		ID:        likeID,
		UserID:    userID,
		PostID:    postID,
		Timestamp: timestamp,
		Active:    true,
	}

	lt.likes[userID] = append(lt.likes[userID], record)
	return likeID
}

// CanUnlike checks if user can unlike a specific post
func (lt *LikeTracker) CanUnlike(postID, userID primitive.ObjectID) *LikeRecord {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	now := time.Now()
	cutoff := now.Add(-3 * time.Minute)

	// Find the most recent active like for this post within window
	userLikes := lt.likes[userID]
	for i := len(userLikes) - 1; i >= 0; i-- {
		record := userLikes[i]
		if record.PostID == postID &&
			record.Active &&
			record.Timestamp.After(cutoff) {
			return &userLikes[i]
		}
	}

	return nil
}

// Unlike marks a specific like record as unliked
func (lt *LikeTracker) Unlike(likeID primitive.ObjectID) bool {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	for userID, records := range lt.likes {
		for i := range records {
			if records[i].ID == likeID && records[i].Active {
				records[i].Active = false
				lt.likes[userID] = records
				return true
			}
		}
	}

	return false
}

// GetPostLikeCount returns active like count for a post
func (lt *LikeTracker) GetPostLikeCount(postID primitive.ObjectID) int {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	count := 0
	for _, records := range lt.likes {
		for _, record := range records {
			if record.PostID == postID && record.Active {
				count++
			}
		}
	}

	return count
}

// GetLikedUsers returns user IDs who liked a post
func (lt *LikeTracker) GetLikedUsers(postID primitive.ObjectID) []string {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	var userIDs []string
	seen := make(map[primitive.ObjectID]bool)

	for _, records := range lt.likes {
		for _, record := range records {
			if record.PostID == postID && record.Active && !seen[record.UserID] {
				userIDs = append(userIDs, record.UserID.Hex())
				seen[record.UserID] = true
			}
		}
	}

	return userIDs
}

// HasUserLikedPost checks if user currently has active like on post
func (lt *LikeTracker) HasUserLikedPost(postID, userID primitive.ObjectID) bool {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	for _, record := range lt.likes[userID] {
		if record.PostID == postID && record.Active {
			return true
		}
	}

	return false
}

// Cleanup removes records older than 3 minutes
func (lt *LikeTracker) Cleanup() {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	cutoff := time.Now().Add(-3 * time.Minute)

	for userID, records := range lt.likes {
		var recent []LikeRecord
		for _, record := range records {
			if record.Timestamp.After(cutoff) {
				recent = append(recent, record)
			}
		}

		if len(recent) == 0 {
			delete(lt.likes, userID)
		} else {
			lt.likes[userID] = recent
		}
	}
}
