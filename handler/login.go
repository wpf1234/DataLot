package handler

import (
	"crypto/tls"
	"datalot/base"
	"datalot/models"
	"datalot/utils"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (g *Gin) Login(c *gin.Context) {
	var id int
	var username, interest, pwd, head string
	//var head interface{}
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

	db := base.DB.Raw("select id,username,password,head,interest from user where phone=?", login.Phone)
	err = db.Row().Scan(&id, &username, &pwd, &head, &interest)
	if err != nil {
		log.Error("查询失败: ", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "请求失败!",
		})
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
	password := utils.StrMd5(login.Password)
	if password != pwd {
		log.Warn("登录密码不正确!!!")
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    nil,
			"message": "密码不正确!",
		})
		return
	}

	//if head == nil {
	//	head = ""
	//}
	//var avatar string
	//_=json.Unmarshal(head.([]uint8),&avatar)
	//sid := base.MapConf.ServiceId
	userId := strconv.Itoa(id)
	userSig, _ := utils.GenSig(int(base.IMC.Appid), base.IMC.Key, userId, expire)
	claims := &utils.MyClaims{
		Id:        id,
		Username:  username,
		Password:  password,
		Mobile:    login.Phone,
		Head:      head,
		Meid:      login.Meid,
		PhoneDesc: login.Desc,
		//ServiceId: sid,
		//ServiceId: sid,
		UserSig: userSig,
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

	// 获取 MapBox 用户的 token
	//mapToken, err := createToken(username)
	//if err != nil {
	//	log.Error("获取 MapBox 的 token 失败: ", err)
	//	c.JSON(http.StatusOK, gin.H{
	//		"code":    http.StatusInternalServerError,
	//		"data":    err,
	//		"message": "MapBox token error!",
	//	})
	//	return
	//}
	//fmt.Println("map box: ",mapToken)

	res := models.LoginRes{
		User: models.User{
			Id:        id,
			Username:  username,
			Phone:     login.Phone,
			Meid:      login.Meid,
			PhoneDesc: login.Desc,
			Head:      head,
			Interest:  interests,
			UserSig:   userSig,
		},
		Token:    token,
		//MapToken: mapToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    res,
		"message": "登录成功!",
	})

	t := time.Now().UnixNano() / 10e5
	db = base.DB.Exec("update user set meid=?,phone_desc=?,login_time=?,user_sig=?,map_token=? where id=?",
		login.Meid, login.Desc, t, userSig, "mapToken" ,id)

	fmt.Println("Update: ", db.RowsAffected)

}

// 微信登录
func (g *Gin) WXLogin(c *gin.Context) {

}

// 验证码登录
func (g *Gin) AuthCodeLogin(c *gin.Context) {

}

// 为每一个用户生成一个 mapbox 的 public token
func createToken(user string) (string, error) {
	var mapData models.MapToken
	var mapToken string
	uri := base.MapBox.Url + base.MapBox.Sk
	mapData.Note = user+" map"
	mapData.Scopes = []string{
		"styles:tiles",
		"styles:read",
		"fonts:read",
		"datasets:read",
		"vision:read",
	}
	mapData.AllowedUrls = []string{
		"*",
	}
	data, _ := json.Marshal(mapData)

	req, err := http.NewRequest("POST", uri, strings.NewReader(string(data)))
	if err != nil {
		log.Error("创建请求失败: ", err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "keep-alive")

	c := &http.Client{
		Transport:     &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, // 跳过https认证
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	response, err := c.Do(req)
	if err != nil {
		log.Error("请求失败: ", err)
		return "", err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	res := make(map[string]interface{})
	_ = json.Unmarshal(body, &res)
	for k, v := range res {
		if k == "token" {
			mapToken = v.(string)
		}
	}
	return mapToken, nil
}
