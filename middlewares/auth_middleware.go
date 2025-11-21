package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"icctv-http-service/services"
)

type contextKey string

const adminClaimsKey contextKey = "adminClaims"

// AuthMiddleware JWT 鉴权
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware 构造函数
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// RequireAdmin 鉴权
func (m *AuthMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			respondUnauthorized(w, "missing bearer token")
			return
		}

		claims, err := m.authService.ValidateToken(token)
		if err != nil {
			respondUnauthorized(w, err.Error())
			return
		}

		ctx := r.Context()
		ctx = contextWithClaims(ctx, claims)
		next(w, r.WithContext(ctx))
	}
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Fields(header)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func contextWithClaims(ctx context.Context, claims *services.AdminClaims) context.Context {
	return context.WithValue(ctx, adminClaimsKey, claims)
}

// AdminClaimsFromContext 读取管理员信息
func AdminClaimsFromContext(ctx context.Context) *services.AdminClaims {
	if val, ok := ctx.Value(adminClaimsKey).(*services.AdminClaims); ok {
		return val
	}
	return nil
}
