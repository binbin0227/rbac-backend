package controllers

import (
	"context"
	"fmt"
	"net/http"
	"rbac-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ctx = context.Background()

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 1. 查询账号密码是否正确
	var user models.User
	if err := models.DB.Preload("Roles.Permissions").Where("name = ? AND password = ?", req.Name, req.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 2. 账号密码正确则生成随机Token
	token := uuid.New().String()

	// 3. Redis: 把 Token 和 UserID 存起来 (有效期 2 小时)
	models.RDB.Set(ctx, token, user.ID, 2*time.Hour)

	// 4. Redis: 把他的所有“权限钥匙”存进 Redis 集合 里
	permKey := fmt.Sprintf("user_perms:%d", user.ID) // 例如 user_perms:1
	models.RDB.Del(ctx, permKey)                     // 先清空旧的缓存
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			// 把 "GET:/api/v1/users" 这种格式存进去
			permStr := perm.Method + ":" + perm.Path
			models.RDB.SAdd(ctx, permKey, permStr)
		}
	}
	models.RDB.Expire(ctx, permKey, 2*time.Hour) // 权限集合有效期 2 小时
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}
