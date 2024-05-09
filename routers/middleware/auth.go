package middleware

import (
	"net/http"
	"strings"

	"gin-boilerplate/config"
	"gin-boilerplate/helpers"
	"gin-boilerplate/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func IsStringInMap(str string, m map[string]models.RoleID) bool {
	_, exists := m[str]
	return exists
}

// 用户角色验证中间件，满足给定条件的才可以修改个人信息
func UserRoleAuthMiddleware(allowed_roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 所有role必须在model中给出
		for _, role := range allowed_roles {
			if !IsStringInMap(role, models.RoleStrToEnumMap) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "角色验证中间件被指派了无效的角色"})
				c.Abort()
				return
			}
		}

		// 从请求头中获取令牌
		tokenString := c.GetHeader("Authorization")
		// 验证令牌
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求未携带令牌"})
			c.Abort()
			return
		}

		// 解析令牌
		token, err := jwt.ParseWithClaims(tokenString, &helpers.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret()), nil
		})

		// 检查错误
		if err != nil {
			// 如果错误是过期错误，则返回特定的过期响应
			if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已过期"})
				c.Abort()
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
				c.Abort()
				return
			}
		}

		// 验证令牌是否有效
		allowed_roles_str := strings.Join(allowed_roles, ",")
		if claims, ok := token.Claims.(*helpers.Claims); ok && token.Valid {
			// 检查用户角色是否尚未分配角色
			if !strings.Contains(allowed_roles_str, claims.UserRole) {
				c.JSON(http.StatusForbidden, gin.H{"error": "当前角色" + claims.UserRole + "无权访问该组路由, 当前路由组允许的角色为: " + allowed_roles_str})
				c.Abort()
				return
			}
			// 继续处理请求
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
		}
	}
}
