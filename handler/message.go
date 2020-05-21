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

// 别人点赞和评论的消息
func (g *Gin) GetMessage(c *gin.Context) {
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
	myName := claims.(*utils.MyClaims).Username
	myHead := claims.(*utils.MyClaims).Head

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	var list []models.Message
	db := base.DB.Raw(fmt.Sprintf(`select * from 
(select id,d_id,operator,operate,tm,is_read from record where u_id=%d ) as a 
left join 
(select id,content,picture,create_at,f_num from dynamic where report=0) as b 
on a.d_id=b.id 
order by a.tm desc,a.is_read asc 
limit %d,%d`, userId, (page-1)*size, size))

	rows, err := db.Rows()

	if err != nil {
		log.Error("查询失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "查询失败!",
		})
		return
	}

	for rows.Next() {
		var msg models.Message
		var d models.Dynamic
		var dId, like, read int
		var p string
		_ = rows.Scan(&msg.MsgId, &dId, &msg.UserId, &msg.Operate, &msg.Tm, &read,
			&d.Id, &d.Content, &p, &d.Tm, &d.Favorite)
		db = base.DB.Raw("select username,head from user where id=?", msg.UserId)
		_ = db.Row().Scan(&msg.Username, &msg.Head)
		msg.Head = preview + msg.Head
		d.Picture = strings.Split(p, ",")
		d.UserId = userId
		d.Username = myName
		d.Head = preview + myHead
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
		msg.Dynamic = d

		list = append(list, msg)
	}
	_ = rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    list,
		"message": "获取消息成功!",
	})
}

// 改变消息已读/未读状态
func (g *Gin) ChangeState(c *gin.Context) {
	id := c.Query("id")
	db := base.DB.Exec("update record set is_read=? where id=?", 1, id)
	err := db.Error
	if err != nil {
		log.Error("更新失败: ", err)
		return
	}
	fmt.Println("Update: ", db.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    nil,
		"message": "更新成功!",
	})
}

// 我喜欢的
func (g *Gin) GetMyFavorite(c *gin.Context) {
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

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	var list []models.Dynamic
	db := base.DB.Raw(fmt.Sprintf(`select * from 
(select d_id from favorite where u_id=%d and is_like=%d) as a 
left join 
(select id,u_id,content,picture,create_at,f_num from dynamic where report=0) as b 
on a.d_id=b.id 
order by b.create_at desc 
limit %d,%d`, userId, 1, (page-1)*size, size))

	rows, err := db.Rows()
	if err != nil {
		log.Error("查询失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "查询失败!",
		})
		return
	}

	for rows.Next() {
		var d models.Dynamic
		var dId int
		var p string
		_ = rows.Scan(&dId, &d.Id, &d.UserId, &d.Content, &p, &d.Tm, &d.Favorite)
		db = base.DB.Raw("select username,head from user where id=?", d.UserId)
		_ = db.Row().Scan(&d.Username, &d.Head)
		d.Head = preview + d.Head
		d.Picture = strings.Split(p, ",")
		db = base.DB.Raw("select count(1) from comment where d_id=?", d.Id)
		_ = db.Row().Scan(&d.Comment)
		d.Like = true
		comm, _ := comment(d.Id)
		d.CommList = comm

		list = append(list, d)
	}
	defer rows.Close()

	fmt.Println("Favorite list: ",len(list))
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    list,
		"message": "获取消息成功!",
	})
}

// 我评论的
func (g *Gin) GetMyComm(c *gin.Context) {
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
	username := claims.(*utils.MyClaims).Username
	myHead := claims.(*utils.MyClaims).Head

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	var list []models.MyComm

	db := base.DB.Raw(fmt.Sprintf(`select * from 
(select id,d_id,comment,comm_tm from comment where u_id=%d) as a 
left join 
(select id,u_id,content,picture,create_at,f_num from dynamic where report=0) as b 
on a.d_id=b.id 
order by a.comm_tm desc 
limit %d,%d`, userId, (page-1)*size, size))

	rows, err := db.Rows()
	if err != nil {
		log.Error("查询失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "查询失败!",
		})
		return
	}

	for rows.Next() {
		var c models.MyComm
		var d models.Dynamic
		var dId, like int
		var p string
		_ = rows.Scan(&c.Id, &dId, &c.Comment, &c.Tm,
			&d.Id, &d.UserId, &d.Content, &p, &d.Tm, &d.Favorite)
		d.Picture = strings.Split(p, ",")
		db = base.DB.Raw("select username,head from user where id=?", d.UserId)
		_ = db.Row().Scan(&d.Username, &d.Head)
		d.Head = preview + d.Head
		db = base.DB.Raw("select count(1) from comment where d_id=?", d.Id)
		_ = db.Row().Scan(&d.Comment)
		d.Like = true
		comm, _ := comment(d.Id)
		d.CommList = comm
		db = base.DB.Raw("select is_like from favorite where u_id=? and d_id=?",
			userId, d.Id)
		_ = db.Row().Scan(&like)
		if like == 1 {
			d.Like = true
		} else {
			d.Like = false
		}
		c.Dynamic = d
		c.Username = username
		c.Head = preview + myHead
		list = append(list, c)
	}
	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    list,
		"message": "获取消息成功!",
	})
}
