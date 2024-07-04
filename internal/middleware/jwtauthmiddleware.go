package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 自定义上下文键类型
type contextKey string

const ClaimsKey = contextKey("claims")

type MyClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JwtAuthMiddleware struct {
	secret string
}

func NewJwtAuthMiddleware(secret string) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		secret: secret,
	}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		const BearerSchema = "Bearer "
		if !strings.HasPrefix(authHeader, BearerSchema) {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		tokenString := authHeader[len(BearerSchema):]

		// 解析和验证 JWT token
		claims := &MyClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 确保 token 的方法符合 "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(m.secret), nil
		})

		if err != nil || !token.Valid {
			logx.Errorf("Invalid token: %v", err)
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		// 将解析后的 claims 信息存储到上下文中，便于后续使用
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}
