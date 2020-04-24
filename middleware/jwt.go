package middleware

import (
	. "datalot/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			log.Error("Token is null!!")

			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"data":    data,
				"message": "未授权!",
			})
			c.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil {
			log.Error("Token解析失败!")
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusInternalServerError,
				"data":    data,
				"message": "Token解析失败!",
			})
			c.Abort()
			return
		}
		// Token过期
		if time.Now().Unix() > claims.ExpiresAt{
			log.Error("Token已过期!")
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusInternalServerError,
				"data":    data,
				"message": "Token已过期!",
			})
			c.Abort()
			return
		}
		// 继续交由下一路由处理，并将解析出的信息传递下去
		c.Set("claims", claims)
		//c.Next()
	}

}
