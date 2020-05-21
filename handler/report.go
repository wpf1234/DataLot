package handler

import (
	"datalot/base"
	"datalot/models"
	"datalot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (g *Gin) Report(c *gin.Context){
	var request models.Report
	err:=c.BindJSON(&request)
	if err!=nil{
		log.Error("请求失败: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusBadRequest,
			"data": err,
			"message": "请求失败!",
		})
		return
	}

	reason:=utils.Slice2Str(request.Reason)
	db:=base.DB.Exec("update dynamic set report=?,reason=?",1,reason)
	fmt.Println("[Report]Update: ",db.RowsAffected)

	c.JSON(http.StatusOK,gin.H{
		"code": http.StatusOK,
		"data": db,
		"message": "已举报!",
	})
}
