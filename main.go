package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"rbac-backend/models"
	"rbac-backend/router"
)

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/rbac_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	models.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("连接数据库成功")
	
	// 建表
	err = models.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})
	if err != nil {
		panic(err)
	}
	fmt.Println("数据库表结构同步完成")

	{
	// ==================== 【新增：创世种子脚本】 ====================
	var count int64
	models.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		fmt.Println("🌱 检测到空数据库，正在植入上帝角色...")
		perms := []models.Permission{
			{Name: "获取用户", Path: "/api/v1/users", Method: "GET"},
			{Name: "创建用户", Path: "/api/v1/users", Method: "POST"},
			{Name: "删除用户", Path: "/api/v1/users/:id", Method: "DELETE"},
			{Name: "获取角色", Path: "/api/v1/roles", Method: "GET"},
			{Name: "创建角色", Path: "/api/v1/roles", Method: "POST"},
			{Name: "删除角色", Path: "/api/v1/roles/:id", Method: "DELETE"},
			{Name: "获取权限", Path: "/api/v1/perms", Method: "GET"},
			{Name: "创建权限", Path: "/api/v1/perms", Method: "POST"},
			{Name: "删除权限", Path: "/api/v1/perms/:id", Method: "DELETE"},
			{Name: "分配角色", Path: "/api/v1/users/:id/roles", Method: "POST"},
			{Name: "分配权限", Path: "/api/v1/roles/:id/perms", Method: "POST"},
		}
		models.DB.Create(&perms)

		// 2. 铸造至高无上的工牌
		adminRole := models.Role{Name: "超级管理员", Description: "创世神，拥有所有权限"}
		models.DB.Create(&adminRole)

		// 3. 把这 11 把钥匙全部拴在这个工牌上
		models.DB.Model(&adminRole).Association("Permissions").Append(perms)

		// 4. 创造上帝本帝 (强制它的 ID 会是 1)
		godUser := models.User{Name: "Admin_Root"}
		models.DB.Create(&godUser)

		// 5. 把工牌挂在上帝脖子上
		models.DB.Model(&godUser).Association("Roles").Append([]models.Role{adminRole})

		fmt.Println("创世完成！上帝账号 Admin_Root (ID:1) 已就绪！")
	}
	// ==============================================================
	}

	// 初始化路由并启动服务
	r := router.SetupRouter()
	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("RBAC 后端服务器启动在 http://127.0.0.1:8080")
}
