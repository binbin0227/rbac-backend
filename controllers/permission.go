package controllers

import (
	"net/http"
	"rbac-backend/models"

	"github.com/gin-gonic/gin"
)

func GetPermissions(c *gin.Context) {
	var permissions []models.Permission
	models.DB.Find(&permissions)
	c.JSON(http.StatusOK, permissions)
}

func CreatePermission(c *gin.Context) {
	var permission models.Permission
	err := c.ShouldBindJSON(&permission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析JSON失败"})
		return
	}
	models.DB.Create(&permission)
	c.JSON(http.StatusOK, permission)
}
func DeletePermission(c *gin.Context) {
	id := c.Param("id")
	models.DB.Where("id = ?", id).Delete(&models.Permission{})
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
