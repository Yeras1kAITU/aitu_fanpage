package service

import (
	"errors"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/config"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/dto"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/middleware"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password must be at least 8 characters long")
)

type AuthService struct {
	userRepo repository.UserRepository
	authMid  *middleware.AuthMiddleware
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		authMid:  middleware.NewAuthMiddleware(cfg, userRepo),
		cfg:      cfg,
	}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*models.User, error) {
	if err := s.validateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := s.validatePassword(req.Password); err != nil {
		return nil, err
	}

	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	user, err := models.NewUser(req.Email, req.Password, req.DisplayName, models.RoleStudent)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if !user.ValidatePassword(req.Password) {
		return "", nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", nil, errors.New("account is deactivated")
	}

	token, err := s.authMid.GenerateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) GetCurrentUser(userID primitive.ObjectID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) UpdateProfile(userID primitive.ObjectID, req dto.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) ChangePassword(userID primitive.ObjectID, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if !user.ValidatePassword(req.CurrentPassword) {
		return errors.New("current password is incorrect")
	}

	if err := s.validatePassword(req.NewPassword); err != nil {
		return err
	}

	if err := user.UpdatePassword(req.NewPassword); err != nil {
		return err
	}

	return s.userRepo.Update(user)
}

func (s *AuthService) validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func (s *AuthService) validatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}
	return nil
}

func (s *AuthService) GetTokenAuth() *middleware.AuthMiddleware {
	return s.authMid
}
