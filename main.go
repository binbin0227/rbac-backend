package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"rbac-backend/models"
	"rbac-backend/router"
)

func main() {
	// 连接Redis
	models.RDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})
	fmt.Println("连接 Redis 成功")

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
		// 创世种子脚本
		var count int64
		models.DB.Model(&models.User{}).Count(&count)
		if count == 0 {
			fmt.Println("检测到空数据库为空，正在植入电商业务预设数据...")
			
			perms := []models.Permission{
				{Name: "查看数据看板", Path: "/api/v1/dashboard", Method: "GET"},
				{Name: "查看商品列表", Path: "/api/v1/products", Method: "GET"},
				{Name: "上架新商品", Path: "/api/v1/products", Method: "POST"},
				{Name: "强制下架商品", Path: "/api/v1/products/:id", Method: "DELETE"},
				{Name: "查看全球订单", Path: "/api/v1/orders", Method: "GET"},
				{Name: "变更订单状态", Path: "/api/v1/orders/:id", Method: "PUT"},
				{Name: "查看财务报表", Path: "/api/v1/finance", Method: "GET"},
				{Name: "审批退款打款", Path: "/api/v1/refunds", Method: "POST"},
				{Name: "分配员工岗位", Path: "/api/v1/users/:id/roles", Method: "POST"},
				{Name: "调整岗位权限", Path: "/api/v1/roles/:id/perms", Method: "POST"},
			}
			models.DB.Create(&perms)

			roles := []models.Role{
				{Name: "超级管理员", Description: "电商平台总架构，拥有所有权限"},
				{Name: "商品运营总监", Description: "负责商品上下架管理，无权碰订单和财务"},
				{Name: "一线客服专员", Description: "负责查看订单和发货，无权碰商品"},
				{Name: "财务审计员", Description: "负责报表和退款审核"},
			}
			models.DB.Create(&roles) 

			models.DB.Model(&roles[0]).Association("Permissions").Append(perms) 
			models.DB.Model(&roles[1]).Association("Permissions").Append([]models.Permission{perms[0], perms[1], perms[2], perms[3]})
			models.DB.Model(&roles[2]).Association("Permissions").Append([]models.Permission{perms[0], perms[4], perms[5]})
			models.DB.Model(&roles[3]).Association("Permissions").Append([]models.Permission{perms[0], perms[6], perms[7]})

			users := []models.User{
				{Name: "Admin_Root", Password: "admin"},
				{Name: "Ops_Alice", Password: "123"},
				{Name: "CS_Bob", Password: "123"},
				{Name: "Finance_Charlie", Password: "123"},
			}
			models.DB.Create(&users)

			models.DB.Model(&users[0]).Association("Roles").Append([]models.Role{roles[0]})
			models.DB.Model(&users[1]).Association("Roles").Append([]models.Role{roles[1]})
			models.DB.Model(&users[2]).Association("Roles").Append([]models.Role{roles[2]})
			models.DB.Model(&users[3]).Association("Roles").Append([]models.Role{roles[3]})

			fmt.Println("4位电商员工账号创建完成！")
		}
	}

	// 初始化路由并启动服务
	r := router.SetupRouter()
	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("RBAC 后端服务器启动在 http://127.0.0.1:8080")
}
