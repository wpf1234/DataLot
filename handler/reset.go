package handler

import (
	"datalot/base"
	"datalot/models"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (g *Gin) ResetPwd(c *gin.Context) {
	var reset models.Reset
	err := c.BindJSON(&reset)
	if err != nil {
		log.Error("获取请求参数失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    err,
			"message": "请求失败！",
		})
		return
	}

	db := base.DB.Exec("update user set password=? where phone=?",
		reset.NewPassword, reset.Phone)
	if err := db.Error; err != nil {
		log.Error("更新失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "重置密码失败！",
		})
		return
	}
	fmt.Println("Update: ", db.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    reset,
		"message": "重置密码成功！",
	})
}
