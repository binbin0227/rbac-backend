package middleware

import (
	"context"
	"fmt"
	"net/http"
	"rbac-backend/models"

	"github.com/gin-gonic/gin"
)
var ctx = context.Background()

func RBACGuard() gin.HandlerFunc{
	return func(c *gin.Context) {
		// 1. 查身份：从 Header 里找 Authorization (前端传 Token)
		token:= c.GetHeader("Authorization")
		if token == ""{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"未登录或缺少 Token"})
			return 
		}

		// 2. 去 Redis 查这个 Token 对应哪个 UserID
		userId,err:=models.RDB.Get(ctx,token).Result()
		if err!=nil{
			c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"登录已过期，请重新登录"})
			return 
		}

		// 3. 看目的
		requestPath:=c.FullPath()
		requestMethod:=c.Request.Method
		requiredPerm:=requestMethod+":"+requestPath // 比如 "DELETE:/api/v1/users/:id"

		// 4. 去 Redis 的 Set 里查询他有没有权限
		permKey:=fmt.Sprintf("user_perms:%s",userId)
		hasPerm,_:=models.RDB.SIsMember(ctx,permKey,requiredPerm).Result()

		if hasPerm{
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden,gin.H{
				"error":"没有 ["+requiredPerm+"] 的权限",
			})
			return 
		}
	}
}