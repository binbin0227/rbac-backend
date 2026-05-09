package controllers

import (
	"net/http"
	"rbac-backend/models"

	"github.com/gin-gonic/gin"
)

type AssignPermRequest struct {
	PermIDs []uint `json:"perm_ids"`
}

func GetRoles(c *gin.Context) {
	var roles []models.Role
	models.DB.Preload("Permissions").Find(&roles)
	c.JSON(http.StatusOK, roles)
}

func CreateRole(c *gin.Context) {
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析JSON失败"})
		return
	}
	models.DB.Create(&role)
	c.JSON(http.StatusOK, role)
}

func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	models.DB.Where("id = ?", id).Delete(&models.Role{})
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 分配权限给角色 (POST /api/v1/roles/:id/perms)
func AssignPermsToRole(c *gin.Context) {
	roleId := c.Param("id")
	var req AssignPermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析JSON失败"})
		return
	}
	var role models.Role
	if err := models.DB.First(&role, roleId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到该角色"})
		return
	}

	var perms []models.Permission
	models.DB.Where("id IN ?", req.PermIDs).Find(&perms)

	// 把找出来的权限，分配给这个角色
	models.DB.Model(&role).Association("Permissions").Replace(perms)

	c.JSON(http.StatusOK, gin.H{"message": "权限分配成功！"})
}
