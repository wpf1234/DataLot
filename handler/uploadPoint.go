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
	"strconv"
	"strings"
)

func (g *Gin) UploadPoint(c *gin.Context)  {
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

	var location models.Location
	err:=c.BindJSON(&location)
	if err!=nil{
		log.Error("请求失败: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusBadRequest,
			"data": err,
			"message": "请求失败!",
		})
		return
	}

	data,err:=json.Marshal(location)
	if err!=nil{
		log.Error("数据转换失败: ",err)
		c.JSON(http.StatusOK,gin.H{
			"code": http.StatusInternalServerError,
			"data": err,
			"message": "数据转换失败!",
		})
		return
	}
	var id int
	var point string
	db:=base.DB.Raw("select id,points from point where u_id=? and begin <=? and end >= ? limit 1" ,
		userId,location.Tm,location.Tm)
	_=db.Row().Scan(&id,&point)
	if point == ""{
		point =string(data)
	}else{
		point = point+","+string(data)
	}
	db=base.DB.Exec("update point set points=? where id=?",point,id)
	fmt.Println("[Point] Update: ",db.RowsAffected)

	c.JSON(http.StatusOK,gin.H{
		"code": http.StatusOK,
		"data": point,
		"message": "上传成功!",
	})

	go getRecommend(id,location.Tm)
}

type rec struct {
	Id int
	Hobby string
	Point string
	Exponent int
}

func getRecommend(id int,t int64){
	// 获取我的坐标点
	var point string
	db:=base.DB.Raw("select points from point where u_id=? and begin <=? and end>=? limit 1",
		id,t,t)
	_=db.Row().Scan(&point)
	myLocal:=strings.Split(point,",")

	var interest,friend string
	db=base.DB.Raw("select interest,friend from user where id=?",id)
	_=db.Row().Scan(&interest,&friend)
	myInter:=strings.Split(interest,",")

	var recommend []rec
	friend = friend +","+strconv.Itoa(id)
	db=base.DB.Raw(fmt.Sprintf("select id,interest from user where id not in(%s)",friend))
	rows,err:=db.Rows()
	if err!=nil{
		log.Error("查询失败: ",err)
		return
	}
	for rows.Next(){
		var uId int
		var uInter string
		var rec rec
		_=rows.Scan(&uId,&uInter)
		uInterest:=strings.Split(uInter,",")
		for _,v:=range myInter{
			for _,vv:=range uInterest{
				if v == vv{
					rec.Id = uId
					rec.Hobby =rec.Hobby + v+","

				}
			}
		}
		rec.Exponent = 1
		rec.Hobby = strings.TrimRight(rec.Hobby,",")
		recommend =append(recommend,rec)
	}

	_=rows.Close()

	for _,v:=range recommend{
		var point string
		var my models.Location
		var other models.Location
		db=base.DB.Raw("select points from point where u_id=? and begin <=? and end>=? limit 1",
			v.Id,t,t)
		_=db.Row().Scan(&point)
		points:=strings.Split(point,",")
		for _,vv:=range myLocal{
			_=json.Unmarshal([]byte(vv),&my)
			for _,val:=range points{
				_=json.Unmarshal([]byte(val),&other)
				if my.Lat == other.Lat && my.Lng == other.Lng{
					v.Point =v.Point + vv + ","
				}
			}
		}
		v.Exponent = 2
		v.Point = strings.TrimRight(v.Point,",")

		var rId int
		var rPoint,rHobby string
		db=base.DB.Raw("select id,points,hobby from recommend where u_id=? and r_id=? ",
			id,v.Id)
		_=db.Row().Scan(&rId,&rPoint,&rHobby)
		if rId == 0{
			db=base.DB.Exec("insert into recommend set u_id=?,r_id=?,points=?,hobby=?,exponent=?",
				id,v.Id,v.Point,v.Hobby,v.Exponent)
			fmt.Println("[rec]Insert: ",db.RowsAffected)
		}else{
			rPoint =rPoint +","+v.Point
			rHobby = rHobby+","+v.Hobby
			db=base.DB.Exec("update into recommend set points=?,hobby=?,exponent=?",
				rPoint,rHobby,v.Exponent)
			fmt.Println("[rec]Insert: ",db.RowsAffected)
		}
	}


}