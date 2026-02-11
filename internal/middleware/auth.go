package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/jwtauth/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Yeras1kAITU/aitu_fanpage/internal/config"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/models"
	"github.com/Yeras1kAITU/aitu_fanpage/internal/repository"
)

type AuthMiddleware struct {
	tokenAuth *jwtauth.JWTAuth
	userRepo  repository.UserRepository
	cfg       *config.Config
}

func NewAuthMiddleware(cfg *config.Config, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenAuth: jwtauth.New("HS256", []byte(cfg.JWT.SecretKey), nil),
		userRepo:  userRepo,
		cfg:       cfg,
	}
}

func (am *AuthMiddleware) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractToken(r)

		// If no token provided, continue without authentication
		if tokenString == "" {
			next.ServeHTTP(w, r)
			return
		}

		token, err := jwtauth.VerifyToken(am.tokenAuth, tokenString)
		if err != nil {
			// If token is invalid, still continue but without user context
			// This allows public
			next.ServeHTTP(w, r)
			return
		}

		claims, err := token.AsMap(context.Background())
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := am.userRepo.FindByID(userID)
		if err != nil || !user.IsActive {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "userID", userID)
		ctx = context.WithValue(ctx, "userRole", string(user.Role))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) RequireRole(role models.UserRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value("userRole")
			if userRole == nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if string(role) != userRole.(string) && userRole.(string) != string(models.RoleAdmin) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (am *AuthMiddleware) RequireAnyRole(roles ...models.UserRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Context().Value("userRole")
			if userRole == nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasRole := false
			currentRole := userRole.(string)

			// Admin can do anything
			if currentRole == string(models.RoleAdmin) {
				next.ServeHTTP(w, r)
				return
			}

			for _, role := range roles {
				if currentRole == string(role) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (am *AuthMiddleware) GenerateToken(user *models.User) (string, error) {
	claims := map[string]interface{}{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"role":    string(user.Role),
		"name":    user.DisplayName,
	}

	_, tokenString, err := am.tokenAuth.Encode(claims)
	return tokenString, err
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return r.URL.Query().Get("token")
	}

	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}

	return bearerToken
}

func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value("user").(*models.User)
	return user, ok
}

func GetUserIDFromContext(ctx context.Context) (primitive.ObjectID, bool) {
	userID, ok := ctx.Value("userID").(primitive.ObjectID)
	return userID, ok
}
