package handler

import (
	"datalot/base"
	"datalot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (g *Gin) AddFriend(c *gin.Context){
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

	// 添加的好友ID
	id:=c.Query("id")

	// 查询我的朋友列表
	db:=base.DB.Raw("select friend from user where id=?",userId)
	var f string
	_=db.Row().Scan(&f)

	friends:=strings.Split(f,",")

	friends=append(friends,id)

	friend:=utils.Slice2Str(friends)

	db=base.DB.Exec("update user set friend=? where id=?",friend,userId)
	fmt.Println("Update: ",db.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    friends,
		"message": "添加成功!",
	})
}

func (g *Gin) DelFriend(c *gin.Context)  {
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

	// 添加的好友ID
	id:=c.Query("id")

	fmt.Println("ID: ",id)

	db:=base.DB.Raw("select friend from user where id=?",userId)
	var f string
	_=db.Row().Scan(&f)

	friends:=strings.Split(f,",")

	for k,v:=range friends{
		if v == id{
			friends=append(friends[:k],friends[k+1:]...)
		}
	}

	friend:=utils.Slice2Str(friends)

	db=base.DB.Exec("update user set friend=? where id=?",friend,userId)
	fmt.Println("Update: ",db.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    friends,
		"message": "删除成功!",
	})
}