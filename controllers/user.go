package controllers

import (
	"net/http"
	"rbac-backend/models"

	"github.com/gin-gonic/gin"
)

type AssignRoleRequest struct {
	RoleIDs []uint `json:"role_ids"`
}

// 获取所有用户 (GET /api/v1/users)
func GetUsers(c *gin.Context) {
	var users []models.User
	models.DB.Preload("Roles").Find(&users)
	c.JSON(http.StatusOK, users)
}

// 创建新用户 (POST /api/v1/users)
func CreateUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析JSON失败"})
		return
	}
	models.DB.Create(&user)
	c.JSON(http.StatusOK, user)
}

// 删除用户 (DELETE /api/v1/users/:id)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	models.DB.Where("id = ?", id).Delete(&models.User{})
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 分配角色给用户 (POST /api/v1/users/:id/roles)
func AssignRolesToUser(c *gin.Context) {
	userId := c.Param("id")
	var req AssignRoleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析JSON失败"})
		return
	}
	// 查找用户
	var user models.User
	if err := models.DB.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到该用户"})
		return
	}
	// 把前端传过来的 ID 数组，去查出对应的角色实体
	var roles []models.Role
	models.DB.Where("id in ?", req.RoleIDs).Find(&roles)
	// 把找出来的 roles 数组，强行替换给这个用户的 Roles 属性
	models.DB.Model(&user).Association("Roles").Replace(roles)
	c.JSON(http.StatusOK, gin.H{"message": "角色分配成功！"})
}
