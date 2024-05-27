package middleware_auth_org

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// ApiIsEnabled 接口是否启用
func ApiIsEnabled() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("===============:", context.Request.RequestURI)
		context.Next()
	}
}
