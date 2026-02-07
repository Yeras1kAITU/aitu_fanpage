package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleStudent   UserRole = "student"
	RoleAdmin     UserRole = "admin"
	RoleAlumni    UserRole = "alumni"
	RoleModerator UserRole = "moderator" // New role for content moderation
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	DisplayName  string             `bson:"display_name" json:"display_name"`
	Role         UserRole           `bson:"role" json:"role"`
	ProfileImage string             `bson:"profile_image,omitempty" json:"profile_image,omitempty"`
	Bio          string             `bson:"bio,omitempty" json:"bio,omitempty"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	LastLoginAt  time.Time          `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
	PostCount    int                `bson:"post_count" json:"post_count"`
	LikeCount    int                `bson:"like_count" json:"like_count"`
	CommentCount int                `bson:"comment_count" json:"comment_count"`
}

func NewUser(email, password, displayName string, role UserRole) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		DisplayName:  displayName,
		Role:         role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastLoginAt:  now,
		PostCount:    0,
		LikeCount:    0,
		CommentCount: 0,
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) UpdatePassword(newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) UpdateLastLogin() {
	u.LastLoginAt = time.Now()
}

// Role check methods
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) IsAlumni() bool {
	return u.Role == RoleAlumni
}

func (u *User) IsModerator() bool {
	return u.Role == RoleModerator
}

func (u *User) CanManageUsers() bool {
	return u.IsAdmin() || u.IsModerator()
}

func (u *User) CanManagePosts() bool {
	return u.IsAdmin() || u.IsModerator()
}

func (u *User) CanManageComments() bool {
	return u.IsAdmin() || u.IsModerator()
}

func (u *User) CanCreatePost() bool {
	return u.IsActive && (u.IsAdmin() || u.IsStudent() || u.IsAlumni() || u.IsModerator())
}

func (u *User) CanEditPost(postAuthorID primitive.ObjectID) bool {
	if u.IsAdmin() || u.IsModerator() {
		return true
	}
	return u.ID == postAuthorID && u.IsActive
}

func (u *User) CanDeletePost(postAuthorID primitive.ObjectID) bool {
	if u.IsAdmin() || u.IsModerator() {
		return true
	}
	return u.ID == postAuthorID && u.IsActive
}

func (u *User) CanViewAnalytics() bool {
	return u.IsAdmin()
}

func (u *User) IncrementPostCount() {
	u.PostCount++
	u.UpdatedAt = time.Now()
}

func (u *User) IncrementLikeCount() {
	u.LikeCount++
	u.UpdatedAt = time.Now()
}

func (u *User) IncrementCommentCount() {
	u.CommentCount++
	u.UpdatedAt = time.Now()
}
