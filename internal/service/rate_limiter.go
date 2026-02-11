package service

import (
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RateLimiter tracks user likes with 3-per-3-minutes limit
type RateLimiter struct {
	mu     sync.RWMutex
	counts map[primitive.ObjectID][]time.Time
	window time.Duration
	max    int
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		counts: make(map[primitive.ObjectID][]time.Time),
		window: 3 * time.Minute,
		max:    3,
	}
}

// CanLike checks if user can like
func (rl *RateLimiter) CanLike(userID primitive.ObjectID) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return rl.getRecentCount(userID) < rl.max
}

// RecordLike records like timestamp
func (rl *RateLimiter) RecordLike(userID primitive.ObjectID) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.counts[userID] = append(rl.counts[userID], time.Now())
}

// getRecentCount returns number of likes within window
func (rl *RateLimiter) getRecentCount(userID primitive.ObjectID) int {
	now := time.Now()
	cutoff := now.Add(-rl.window)
	count := 0

	for _, ts := range rl.counts[userID] {
		if ts.After(cutoff) {
			count++
		}
	}

	return count
}

// Cleanup removes old timestamps
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-rl.window)

	for userID, timestamps := range rl.counts {
		var recent []time.Time
		for _, ts := range timestamps {
			if ts.After(cutoff) {
				recent = append(recent, ts)
			}
		}

		if len(recent) == 0 {
			delete(rl.counts, userID)
		} else {
			rl.counts[userID] = recent
		}
	}
}
