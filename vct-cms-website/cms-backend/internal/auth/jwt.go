package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ─── CONFIG ───
// JWT Secret đọc từ environment variable, fallback nếu chưa set
var JWTSecret = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "vct-cms-dev-secret-CHANGE-IN-PRODUCTION"
	}
	return []byte(secret)
}()

var TokenExpiry = 24 * time.Hour

// ─── CUSTOM CLAIMS ───
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // "admin" | "editor"
	jwt.RegisteredClaims
}

// GenerateToken tạo JWT Token cho user đã xác thực
func GenerateToken(userID uint, username, role string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vct-cms",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateToken giải mã và xác thực JWT Token
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}
