package auth

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 定义JWT相关错误
var (
	ErrTokenExpired     = errors.New("令牌已过期")
	ErrTokenNotValidYet = errors.New("令牌尚未生效")
	ErrTokenMalformed   = errors.New("令牌格式错误")
	ErrTokenInvalid     = errors.New("令牌无效")
)

// CustomClaims 自定义声明
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTAuth JWT认证
type JWTAuth struct {
	secret string
	expire time.Duration
}

// NewJWTAuth 创建JWT认证
func NewJWTAuth(secret string, expire int) *JWTAuth {
	return &JWTAuth{
		secret: secret,
		expire: time.Duration(expire) * time.Hour, // 配置中的过期时间单位是小时
	}
}

// GenerateToken 生成令牌
func (j *JWTAuth) GenerateToken(userID uint, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "eden-ops",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ParseToken 解析令牌
func (j *JWTAuth) ParseToken(tokenString string) (*CustomClaims, error) {
	log.Printf("解析token: %s", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		log.Printf("解析token失败: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		log.Printf("解析token成功: userID=%d, username=%s", claims.UserID, claims.Username)
		return claims, nil
	}

	log.Printf("无效的token")
	return nil, errors.New("invalid token")
}

// RefreshToken 刷新令牌
func (j *JWTAuth) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.expire))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}
