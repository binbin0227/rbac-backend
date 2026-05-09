package router

import (
	"github.com/gin-gonic/gin"
	"rbac-backend/controllers"
	"rbac-backend/middleware"
)

// Cors 跨域中间件：允许前端网页跨端口访问 Go 接口
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Id")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// SetupRouter 初始化路由配置
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())
	v1 := r.Group("/api/v1")
	v1.Use(middleware.RBACGuard())
	{
		// Users 路由
		v1.GET("/users", controllers.GetUsers)
		v1.POST("/users", controllers.CreateUser)
		v1.DELETE("/users/:id", controllers.DeleteUser)

		// Roles 路由
		v1.GET("/roles", controllers.GetRoles)
		v1.POST("/roles", controllers.CreateRole)
		v1.DELETE("/roles/:id", controllers.DeleteRole)

		// Permissions 路由
		v1.GET("/perms", controllers.GetPermissions)
		v1.POST("/perms", controllers.CreatePermission)
		v1.DELETE("/perms/:id", controllers.DeletePermission)

		// 核心关系绑定路由
		v1.POST("/users/:id/roles", controllers.AssignRolesToUser)
		v1.POST("/roles/:id/perms", controllers.AssignPermsToRole)
	}
	return r
}
