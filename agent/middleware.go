package agent

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"super-markdown-editor-web/model"
)

// 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Request.Cookie("token")
		if err != nil || token.Value == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": -1, "code": 401, "msg": "error", "data": nil})
			c.Abort()
			return
		}
		j := model.NewJWT()
		claims, err := j.ParserToken(token.Value)
		if err != nil || claims == nil {
			c.JSON(http.StatusOK, gin.H{"status": -1, "code": 401, "msg": err.Error(), "data": nil})
			c.Abort()
			return
		}
		c.Set("claims", claims)
	}
}
