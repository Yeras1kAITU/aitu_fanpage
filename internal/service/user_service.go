package service

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/repository"
)

var (
	ErrPermissionDenied     = errors.New("permission denied")
	ErrInvalidRole          = errors.New("invalid role")
	ErrCannotDeactivateSelf = errors.New("cannot deactivate your own account")
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) GetAllUsers(limit, offset int) ([]*models.User, error) {
	return s.userRepo.FindAll(limit, offset)
}

func (s *UserService) UpdateUserRole(adminID, targetUserID primitive.ObjectID, newRole models.UserRole) error {
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return err
	}

	if !admin.CanManageUsers() {
		return ErrPermissionDenied
	}

	if adminID == targetUserID && newRole != models.RoleAdmin {
		return errors.New("cannot change your own role from admin")
	}

	targetUser, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		return ErrUserNotFound
	}

	// Validate role
	switch newRole {
	case models.RoleAdmin, models.RoleStudent, models.RoleAlumni, models.RoleModerator:
		targetUser.Role = newRole
	default:
		return ErrInvalidRole
	}

	return s.userRepo.Update(targetUser)
}

func (s *UserService) DeactivateUser(adminID, targetUserID primitive.ObjectID) error {
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return err
	}

	if !admin.CanManageUsers() {
		return ErrPermissionDenied
	}

	if adminID == targetUserID {
		return ErrCannotDeactivateSelf
	}

	targetUser, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		return ErrUserNotFound
	}

	targetUser.IsActive = false
	return s.userRepo.Update(targetUser)
}

func (s *UserService) ActivateUser(adminID, targetUserID primitive.ObjectID) error {
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return err
	}

	if !admin.CanManageUsers() {
		return ErrPermissionDenied
	}

	targetUser, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		return ErrUserNotFound
	}

	targetUser.IsActive = true
	return s.userRepo.Update(targetUser)
}

func (s *UserService) DeleteUser(adminID, targetUserID primitive.ObjectID) error {
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return err
	}

	if !admin.IsAdmin() {
		return ErrPermissionDenied
	}

	if adminID == targetUserID {
		return errors.New("cannot delete your own account")
	}

	return s.userRepo.Delete(targetUserID)
}

func (s *UserService) GetUserStats(userID primitive.ObjectID) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	stats := map[string]interface{}{
		"user_id":                 user.ID.Hex(),
		"display_name":            user.DisplayName,
		"role":                    string(user.Role),
		"post_count":              user.PostCount,
		"like_count":              user.LikeCount,
		"comment_count":           user.CommentCount,
		"is_active":               user.IsActive,
		"created_at":              user.CreatedAt,
		"last_login_at":           user.LastLoginAt,
		"days_since_registration": int(time.Since(user.CreatedAt).Hours() / 24),
	}

	return stats, nil
}

func (s *UserService) SearchUsers(query string, limit int) ([]*models.User, error) {
	allUsers, err := s.userRepo.FindAll(1000, 0)
	if err != nil {
		return nil, err
	}

	var filteredUsers []*models.User
	query = strings.ToLower(query)

	for _, user := range allUsers {
		if strings.Contains(strings.ToLower(user.Email), query) ||
			strings.Contains(strings.ToLower(user.DisplayName), query) {
			filteredUsers = append(filteredUsers, user)
		}

		if len(filteredUsers) >= limit {
			break
		}
	}

	return filteredUsers, nil
}
