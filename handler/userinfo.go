package handler

import (
	"datalot/base"
	"datalot/models"
	"datalot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func (g *Gin) GetUserInfo(c *gin.Context){
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
	userId := claims.(*utils.MyClaims).Id // user_id

	id:=c.Query("id")
	if id == ""{
		log.Error("请求失败!")
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusBadRequest,
			"data": id,
			"message": "请求失败!",
		})
		return
	}

	var user models.UserInfo
	var dynamic []models.Dynamic
	var friend,sql string
	db:=base.DB.Raw("select username,head,friends from user where id=?",id)
	_=db.Row().Scan(&user.Username,&user.Head,&friend)
	user.Id,_ = strconv.Atoi(id)
	user.Head = preview + user.Head

	myId:=strconv.Itoa(userId)

	if strings.Contains(friend,myId){
		sql=fmt.Sprintf(`select id,u_id,content,picture,create_at,f_num from dynamic 
where u_id=%d and report=0 and authority in (0,1,2)`,user.Id)
	}else{
		sql=fmt.Sprintf(`select id,u_id,content,picture,create_at,f_num from dynamic 
where u_id=%d and report=0 and authority in (0,1)`,user.Id)
	}

	db=base.DB.Raw(sql)
	rows,err:=db.Rows()
	if err!=nil{
		log.Error("查询失败: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusInternalServerError,
			"data": err,
			"message": "查询失败!",
		})
		return
	}

	for rows.Next(){
		var d models.Dynamic
		var p string
		var like int
		_=rows.Scan(&d.Id,&d.UserId,&d.Content,&p,&d.Tm,&d.Favorite)
		d.Picture = strings.Split(p,",")
		d.Username = user.Username
		d.Head = user.Head
		db = base.DB.Raw("select count(1) from comment where d_id=?", d.Id)
		_ = db.Row().Scan(&d.Comment)
		db = base.DB.Raw("select is_like from favorite where u_id=? and d_id=?",
			userId, d.Id)
		_ = db.Row().Scan(&like)
		if like == 1 {
			d.Like = true
		} else {
			d.Like = false
		}
		comm, _ := comment(d.Id)
		d.CommList = comm

		dynamic=append(dynamic,d)
	}
	defer rows.Close()

	user.Dynamic = dynamic

	c.JSON(http.StatusOK,gin.H{
		"code": http.StatusOK,
		"data": user,
		"message": "获取信息成功!",
	})
}
