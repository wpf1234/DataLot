package handler

import (
	"datalot/base"
	"datalot/models"
	"datalot/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

func (g *Gin) Login(c *gin.Context) {
	var id int
	var username, interest, meid, desc string
	var head interface{}
	var login models.Login
	err := c.BindJSON(&login)
	if err != nil {
		log.Error("获取请求结果失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    err,
			"message": "请求失败!",
		})
		return
	}

	db := base.DB.Raw("select id,username,head,interest,meid,phone_desc from user where password=? and phone=?", login.Password, login.Phone)
	err = db.Row().Scan(&id, &username, &head, &interest, &meid, &desc)
	if err != nil {
		log.Error("查询失败: ", err)
		return
	}
	if id == 0 {
		log.Warn("查询用户信息失败,id: ", id)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    nil,
			"message": "该用户未注册，请前往注册!",
		})
		return
	}

	if head == nil {
		head = ""
	}

	sid := base.MapConf.ServiceId
	claims := &utils.MyClaims{
		Id:        id,
		Username:  username,
		Password:  login.Password,
		Mobile:    login.Phone,
		Head:      head.(string),
		Meid:      meid,
		PhoneDesc: desc,
		//ServiceId: sid,
		ServiceId: sid,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                                             // 签名生效时间
			ExpiresAt: time.Now().Add(time.Duration(7*utils.ExpireTime) * time.Hour).Unix(), // 过期时间
		},
	}

	token, err := utils.GetToken(claims)
	if err != nil {
		log.Error("生成Token失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "生成Token失败!",
		})
		return
	}

	// 将兴趣字符串转为切片
	interests := strings.Split(interest, ",")

	//// 将设备添加到我们的轨迹服务中
	//tid, _ := addTerminal(meid, desc)
	//// 添加轨迹
	//trid, _ := addTrack(tid)

	//res := models.LoginRes{
	//	User: models.User{
	//		Id:        id,
	//		Username:  username,
	//		Phone:     login.Phone,
	//		Meid:      meid,
	//		PhoneDesc: desc,
	//		Head:      head,
	//		Interest:  interests,
	//	},
	//	Map: models.MapRes{
	//		ServiceId:  sid,
	//		TerminalId: tid,
	//		TrackId:    trid,
	//	},
	//	Token: token,
	//}

	// 将设备添加到百度地图
	err = addEntity(sid, meid, desc)
	if err != nil {
		log.Error("创建百度地图服务失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "创建百度地图服务失败!",
		})
		return
	}

	res := models.LoginRes{
		User: models.User{
			Id:        id,
			Username:  username,
			Phone:     login.Phone,
			Meid:      meid,
			PhoneDesc: desc,
			Head:      head.(string),
			Interest:  interests,
		},
		Map: models.MapRes{
			ServiceId:  sid,
			EntityName: meid,
		},
		Token: token,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    res,
		"message": "登录成功!",
	})

	t := time.Now().UnixNano() / 10e5
	db = base.DB.Exec("update user set login_time=? where id=?", t, id)
	fmt.Println("Update: ", db.RowsAffected)
}

// 微信登录
func (g *Gin) WXLogin(c *gin.Context){

}