package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"seaguard-admin-backend/models"
	"seaguard-admin-backend/utils"
	"strings"
)

type UserGetter interface {
	GetUserByID(id uint) (*models.User, error)
}

// AuthMiddleware 认证中间件
func AuthMiddleware(userService UserGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证token"})
			c.Abort()
			return
		}

		// 检查Bearer token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token格式"})
			c.Abort()
			return
		}

		// 解析JWT token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			log.Printf("Token解析失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 验证用户是否存在
		user, err := userService.GetUserByID(claims.UserID)
		if err != nil {
			log.Printf("用户验证失败: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户"})
			c.Abort()
			return
		}

		// 检查用户状态
		if user.Status != "active" {
			log.Printf("用户 %d 已被禁用", user.ID)
			c.JSON(http.StatusForbidden, gin.H{"error": "用户已被禁用"})
			c.Abort()
			return
		}

		// 验证用户角色是否匹配
		if user.Role != claims.Role {
			log.Printf("用户角色不匹配: token中为 %s, 数据库中为 %s", claims.Role, user.Role)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户角色验证失败"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", user.ID)
		c.Set("userRole", user.Role)
		c.Next()
	}
}

// AdminRequired 管理员权限验证中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("userRole")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// VolunteerRequired 志愿者权限验证中间件
func VolunteerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("userRole")
		if role != "volunteer" && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要志愿者权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
