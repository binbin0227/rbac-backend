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
	r.POST("/api/v1/login", controllers.Login) // 登录路由
	v1 := r.Group("/api/v1")
	v1.Use(middleware.RBACGuard())
	{
		{// 电商业务空壳接口
		v1.GET("/dashboard", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "成功访问数据看板，今日营业额 $10000"}) })
		
		v1.GET("/products", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "成功获取商品列表：1.无人机 2.机械键盘"}) })
		v1.POST("/products", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "新商品上架成功！"}) })
		v1.DELETE("/products/:id", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "强制下架违规商品成功！"}) })
		
		v1.GET("/orders", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "成功获取订单列表：待发货 50 单"}) })
		v1.PUT("/orders/:id", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "订单状态变更为：已发货"}) })
		
		v1.GET("/finance", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "成功获取财务报表：利润率 35%"}) })
		v1.POST("/refunds", func(c *gin.Context) { c.JSON(200, gin.H{"msg": "财务退款打款成功，资金已原路退回！"}) })
		}

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
