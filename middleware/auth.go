package middleware

import (
	"net/http"
	"rbac-backend/models"

	"github.com/gin-gonic/gin"
)

func RBACGuard() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. 查身份
		// 在HTTP 请求头 Header 寻找一个叫 X-User-Id 的标签
		userId := c.GetHeader("X-User-Id")
		if userId == "" { // AbortWithStatusJSON 阻止继续执行
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header中缺少 X-User-Id"})
			return
		}

		// 2. 看目的
		requestPath := c.FullPath() // 如 /api/v1/users/1
		requestMethod := c.Request.Method // 如 DELETE

		// 3. 翻档案
		var user models.User
		if err := models.DB.Preload("Roles.Permissions").First(&user, userId).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "查无此人"})
			return
		}

		// 4. 对暗号：开始遍历他的所有权限
		hasPermission := false

		for _, role := range user.Roles { // 翻看他的每一个角色
			for _, perm := range role.Permissions { // 翻看每个角色的每个权限
				// 如果路径和方法，跟当前请求完全对得上
				if perm.Path == requestPath && perm.Method == requestMethod {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		// 5. 做决定
		if hasPermission {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "没有 [" + requestMethod + " " + requestPath + "] 的权限",
			})
			return
		}
	}
}
