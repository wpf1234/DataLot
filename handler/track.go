package handler

import (
	"datalot/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (g *Gin) GetTrack(c *gin.Context) {
	claims, ok := c.Get("claims")
	if !ok {
		log.Error("Claims字段不存在!")
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    nil,
			"message": "没有获取到信息!",
		})
		return
	}
	entityName := claims.(*utils.MyClaims).Meid
	serviceId := claims.(*utils.MyClaims).ServiceId

	// 前端传来的时间戳是 ms
	sTm := c.Query("start")
	eTm := c.Query("end")

	start, _ := strconv.ParseInt(sTm, 10, 64)
	end, _ := strconv.ParseInt(eTm, 10, 64)

	start = start / 10e3
	end = end / 10e3

	data, err := track(serviceId, entityName, start, end)
	if err != nil {
		log.Error("获取轨迹失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "获取轨迹失败!",
		})
		return
	}

	res:=make(map[string]interface{})
	_=json.Unmarshal(data,&res)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    res,
		"message": "获取轨迹成功!",
	})
}
