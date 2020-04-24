package handler

import (
	"datalot/base"
	"datalot/models"
	"datalot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (g *Gin) Register(c *gin.Context) {
	var reg models.Register
	var interest string
	err := c.BindJSON(&reg)
	if err != nil {
		log.Error("获取请求结果失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    err,
			"message": "请求失败!",
		})
		return
	}

	// 将数组转换成字符串
	if reg.Interest == nil {
		interest = ""
	} else {
		interest = utils.Slice2Str(reg.Interest)
	}

	if reg.Username == "" {
		reg.Username = "lot_" + reg.Phone
	}
	reg.Password = utils.StrMd5(reg.Password)
	t := time.Now().UnixNano() / 10e5
	db := base.DB.Exec(`insert into user set create_at=?,username=?,password=?,phone=?,interest=?,
		reg_time=?`,
		t, reg.Username, reg.Password, reg.Phone, interest, t)
	err = db.Error
	if err != nil {
		log.Error("新增用户失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "注册失败!",
		})
		return
	}

	fmt.Println("Insert: ", db.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    nil,
		"message": "注册成功!",
	})
}
