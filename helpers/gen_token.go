package helpers

import (
	"time"

	"gin-boilerplate/config"
	"gin-boilerplate/models"

	"github.com/golang-jwt/jwt"
)

// 含有客户信息的JWT声明
type Claims struct {
	UserName string `json:"username"`
	UserRole string `json:"user_role"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(user models.User) (string, string, error) {
	JWT_SECRET, JWT_ALGORITHM, JWT_ACCESS_TOKEN_EXPIRE_MINUTES, JWT_REFRESH_TOKEN_EXPIRE_MINUTES := config.JWTConfiguration()

	// 创建访问令牌
	accessTokenClaims := &Claims{
		UserName: user.UserName,
		UserRole: models.RoleNameMap[user.RoleID],
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(JWT_ACCESS_TOKEN_EXPIRE_MINUTES)).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.GetSigningMethod(JWT_ALGORITHM), accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", "", err
	}

	// 创建刷新令牌
	refreshTokenClaims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(JWT_REFRESH_TOKEN_EXPIRE_MINUTES)).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod(JWT_ALGORITHM), refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
