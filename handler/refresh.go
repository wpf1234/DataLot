package handler

import (
	"datalot/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (g *Gin) Refresh(c *gin.Context) {
	token := c.Request.Header.Get("token")
	claims, err := utils.ParseToken(token)
	if err != nil {
		log.Error("Token解析失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "Token解析失败!",
		})
		return
	}
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(7*utils.ExpireTime) * time.Hour).Unix()
	// 重新生成腾讯云IM 即时通讯的 UserSig
	//newSig,_:=utils.GenSig(int(base.IMC.Appid),base.IMC.Key,claims.Mobile,expire)
	//claims.UserSig=newSig
	newToken, err := utils.GetToken(claims)
	if err != nil {
		log.Error("刷新失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "Token刷新失败!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    newToken,
		"message": "刷新成功!",
	})
}
