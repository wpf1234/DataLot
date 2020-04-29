package handler

import (
	"datalot/base"
	"datalot/models"
	"datalot/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 发表动态
func (g *Gin) Release(c *gin.Context) {
	var content models.Content
	err := c.BindJSON(&content)
	if err != nil {
		log.Error("获取请求结果失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    err,
			"message": "请求失败!",
		})
		return
	}

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

	//	将动态存入数据库
	picture := utils.Slice2Str(content.Picture)

	t := time.Now().UnixNano() / 10e5
	db := base.DB.Exec("insert into dynamic set u_id=?,content=?,picture=?,authority=?,create_at=?",
		userId, content.Text, picture, content.Auth, t)
	err = db.Error
	if err != nil {
		log.Error("插入数据失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "发布动态失败，请重试!",
		})
		return
	}

	fmt.Println("Insert : ", db.RowsAffected)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusInternalServerError,
		"data":    content,
		"message": "发布动态成功!",
	})
}

// 查看动态列表
func (g *Gin) ViewList(c *gin.Context) {
	var list []models.Dynamic
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

	//	第一步，查询所有符合条件的动态
	// 自己的动态
	db := base.DB.Raw(`select id,content,picture,create_at,f_num from dynamic where u_id=?`, userId)
	rows, err := db.Rows()
	if err != nil {
		log.Error("查询动态失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "获取动态失败!",
		})
		return
	}

	for rows.Next() {
		var l models.Dynamic
		var p string
		_ = rows.Scan(&l.Id, &l.Content, &p, &l.Tm, &l.Favorite)
		l.Username = myName
		l.Head = myHead
		l.Picture = strings.Split(p, ",")
		db = base.DB.Raw("select count(1) from comment where d_id=?", l.Id)
		db.Scan(&l.Comment)
		comm,_:=comment(l.Id)
		l.CommList=comm
		l.Active=false
		list = append(list, l)
	}
	// 好友动态
	var friend string
	db = base.DB.Raw("select friends from user where id=?", userId)
	_ = db.Row().Scan(&friend)
	if friend == "" {
		fmt.Println("朋友为空!")
	} else {
		sql := fmt.Sprintf(`select * from 
		(select id,u_id,content,picture,create_at,f_num from dynamic where u_id in (%s)) as a 
		left join 
		(select id,username,head from user) as b 
		on a.u_id=b.id`, friend)
		db = base.DB.Raw(sql)
		rows, err = db.Rows()
		if err != nil {
			log.Error("查询朋友的动态失败: ", err)
			return
		}
		for rows.Next() {
			var l models.Dynamic
			var uId, uuId int
			var p string
			_ = rows.Scan(&l.Id, &uId, &l.Content, &p, &l.Tm, &l.Favorite,
				&uuId, &l.Username, &l.Head)
			l.Picture = strings.Split(p, ",")
			db = base.DB.Raw("select count(1) from comment where d_id=?", l.Id)
			db.Scan(&l.Comment)
			comm,_:=comment(l.Id)
			l.CommList=comm
			l.Active=false
			list = append(list, l)
		}
	}

	// 公开的
	sql := fmt.Sprintf(`select * from 
		(select id,u_id,content,picture,create_at,f_num from dynamic where authority in (0,1)) as a 
		left join 
		(select id,username,head from user) as b 
		on a.u_id=b.id`)
	db = base.DB.Raw(sql)
	rows, err = db.Rows()
	if err != nil {
		log.Error("查询朋友的动态失败: ", err)
		return
	}

	for rows.Next() {
		var l models.Dynamic
		var uId, uuId int
		var p string
		_ = rows.Scan(&l.Id, &uId, &l.Content, &p, &l.Tm, &l.Favorite,
			&uuId, &l.Username, &l.Head)
		l.Picture = strings.Split(p, ",")
		db = base.DB.Raw("select count(1) from comment where d_id=?", l.Id)
		db.Scan(&l.Comment)
		comm,_:=comment(l.Id)
		l.CommList=comm
		l.Active=false
		list = append(list, l)
	}
	_ = rows.Close()

	distinct := utils.Duplicate(list)

	utils.SortBody(distinct, func(p, q *interface{}) bool {
		v := reflect.ValueOf(*p)
		i := v.FieldByName("Tm")
		v = reflect.ValueOf(*q)
		j := v.FieldByName("Tm")
		return i.String() < j.String()
	})

	fmt.Println(distinct)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    distinct,
		"message": "获取动态成功!",
	})
}

// 查看某一个动态的详情
func comment(dynamicId int) ([]models.Comment,error) {

	// 查询该id 动态的信息
	// 获取基本信息
	//sql := fmt.Sprintf(`select * from
	//	(select id,u_id,content,picture,create_at,f_num from dynamic where id=%d) as a
	//	left join
	//	(select id,username,head from user) as b
	//	on a.u_id=b.id`, dynamicId)
	//db := base.DB.Raw(sql)
	//var dynamic models.Dynamic
	//var picture string
	//var uId, uuid int
	//_ = db.Row().Scan(&dynamic.Id, &uId, &dynamic.Content, &picture, &dynamic.Tm, &dynamic.Favorite,
	//	&uuid, &dynamic.Username, &dynamic.Head)
	//dynamic.Picture = strings.Split(picture, ",")

	// 获取评论信息
	//db := base.DB.Raw("select count(1) from comment where d_id=?", dynamicId)
	//_ = db.Row().Scan(&dynamic.Comment)

	var comments []models.Comment
	db := base.DB.Raw(`select * from 
		(select u_id,comment,comm_tm,reply from comment where d_id=?) as a 
		left join 
		(select id,username,head from user) as b 
		on a.u_id=b.id`, dynamicId)
	rows, err := db.Rows()
	if err != nil {
		log.Error("获取评论人信息失败: ", err)
		return nil,err
	}

	for rows.Next() {
		var comm models.Comment
		var reply []models.Reply
		var r=make([]byte,0)
		var uId int
		_ = rows.Scan(&uId, &comm.Context, &comm.Tm, &r,
			&comm.Id, &comm.CommUser, &comm.Head)
		err = json.Unmarshal(r, &reply)
		if err != nil {
			log.Error("数据转换失败: ", err)
			return nil,err
		}
		comm.Reply = reply
		comm.DynamicId = dynamicId

		comments = append(comments, comm)
	}
	_ = rows.Close()

	return comments,nil
}

// 评论某个动态
func (g *Gin) Comment(c *gin.Context) {
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
	username:= claims.(*utils.MyClaims).Username
	myHead := claims.(*utils.MyClaims).Head

	var comm models.WriteComm
	err:=c.BindJSON(&comm)
	if err!=nil{
		log.Error("获取请求参数失败: ",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    err,
			"message": "没有获取到信息!",
		})

		return
	}
	comm.UserId=userId
	//comm.Head=myHead

	// 将数据存入数据库
	tm:=time.Now().UnixNano()/10e5
	db:=base.DB.Exec("insert into comment set d_id=?,u_id=?,username=?,head=?,comment=?,comm_tm=?",
		comm.DynamicId,comm.UserId,username,myHead,comm.Context,tm)
	err=db.Error
	if err!=nil{
		log.Error("插入失败: ",err)
		return
	}
	fmt.Println("Insert : ",db.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    comm,
		"message": "评论成功!",
	})
}

// 答复某条评论
func (g *Gin) Reply(c *gin.Context){
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
	username:= claims.(*utils.MyClaims).Username

	var re models.WriteReply
	err:=c.BindJSON(&re)
	if err!=nil{
		log.Error("获取请求参数失败: ",err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    err,
			"message": "没有获取到信息!",
		})
		return
	}

	// 首先查询该条评论是否有答复
	db:=base.DB.Raw("select reply from comment where id=?",re.CommId)
	var reply=make([]byte,0)
	_=db.Row().Scan(&reply)

	var replays []models.Reply
	if len(reply) == 0{
		// 暂无评论
		var newReply models.Reply
		newReply.Id=userId
		newReply.ReplyUser=username
		newReply.Context=re.Context
		newReply.Tm=re.Tm

		replays=append(replays,newReply)
		reply,err=json.Marshal(replays)
		if err!=nil{
			log.Error("数据转换失败: ",err)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusInternalServerError,
				"data":    err,
				"message": "数据转换失败!",
			})
			return
		}

		db=base.DB.Exec("insert into comment set reply=? where id=?",reply,re.CommId)
		err=db.Error
		if err!=nil{
			log.Error("新增失败: ",err)
			return
		}
		fmt.Println("Insert: ",db.RowsAffected)
	}else{
		// 已经有评论
		_=json.Unmarshal(reply,&replays)

		var newReply models.Reply
		newReply.Id=userId
		newReply.ReplyUser=username
		newReply.Context=re.Context
		newReply.Tm=re.Tm

		replays=append(replays,newReply)

		reply,err=json.Marshal(replays)
		if err!=nil{
			log.Error("数据转换失败: ",err)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusInternalServerError,
				"data":    err,
				"message": "数据转换失败!",
			})
			return
		}

		db=base.DB.Exec("update comment set reply=? where id=?",reply,re.CommId)
		err=db.Error
		if err!=nil{
			log.Error("更新失败: ",err)
			return
		}
		fmt.Println("Update: ",db.RowsAffected)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    re,
		"message": "评论成功!",
	})
}

// 点赞、取消
func (g *Gin) Tags (c *gin.Context){
	id,_:=strconv.Atoi(c.Query("id"))
	like,_:=strconv.Atoi(c.Query("like"))
	if id == 0{
		log.Error("请求参数错误,id=",id)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    nil,
			"message": "请求参数错误!",
		})
		return
	}

	var count int
	db:=base.DB.Raw("select f_num from dynamic where id=?",id)
	_=db.Row().Scan(&count)

	if like == 1{
		count=count+1
	}else{
		count=count-1
	}


	db=base.DB.Exec("update dynamic set f_num=? where id=?",count,id)
	fmt.Println("Update: ",db.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    nil,
		"message": "成功!",
	})
}