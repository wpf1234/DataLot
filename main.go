package main

import (
	"datalot/base"
	"datalot/handler"
	"datalot/middleware"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/web"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	service := web.NewService(
		web.Name("lot.micro.api.v1.lot"),
		web.Version("latest"),
		web.Address(":12345"),
	)
	err := service.Init()
	if err != nil {
		log.Error("服务初始化失败: ", err)
		return
	}

	base.Init()
	go func() {
		for {
			now := time.Now()
			handler.InitPoints()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
	//handler.CreateService()
	// 全局设置环境，debug 为开发环境，线上环境为 gin.ReleaseMode
	gin.SetMode(gin.ReleaseMode)
	// 创建 Restful handler
	g := new(handler.Gin)
	router := gin.Default()
	router.Use(middleware.Cors())

	noAuth := router.Group("/v1/lot")
	noAuth.POST("/login", g.Login)
	noAuth.POST("/register", g.Register)
	//noAuth.DELETE("/friend",g.DelFriend)
	noAuth.POST("/reset", g.ResetPwd)
	//noAuth.POST("/dynamic/upload", g.CyclePicture)
	noAuth.GET("/preview", g.OnlinePreview)
	noAuth.DELETE("/picture", g.DeletePic)  // 删除本地文件夹中的文件
	noAuth.PUT("/msg/state", g.ChangeState) // 更改状态

	auth := router.Group("/v1/lot/auth")
	auth.Use(middleware.JWTAuth())
	{
		// 动态圈
		auth.GET("/dynamic", g.ViewList) // 查看动态列表
		//auth.GET("/dynamic/one", g.ViewOne)          // 查看某个动态的详情
		auth.PUT("/dynamic/tags", g.Tags)            // 点赞
		auth.POST("/dynamic", g.Release)             // 发表动态
		auth.POST("/dynamic/comm", g.Comment)        // 评论动态
		auth.POST("/dynamic/reply", g.Reply)         // 答复评论
		auth.POST("/dynamic/upload", g.CyclePicture) // 上传图片(可上传多张图片)
		auth.GET("/dynamic/mine", g.GetMyDynamic)

		auth.PUT("/friend", g.AddFriend)    // 添加好友
		auth.DELETE("/friend", g.DelFriend) // 删除好友
		auth.GET("/friend", g.FriendList)

		//auth.GET("/track", g.GetTrack) // 获取轨迹

		auth.GET("/refresh", g.Refresh) // 更新Token值

		auth.POST("/upload", g.HeadPortrait) // 上传头像

		auth.POST("/location", g.UploadPoint) // 定位点上传

		auth.GET("/total", g.GetMsgNum) // 获取未读消息数量

		auth.GET("/notice", g.GetMessage) // 别人对于我的动态的点赞和评论
		auth.GET("/my/favorite", g.GetMyFavorite)
		auth.GET("/my/comment", g.GetMyComm)

		auth.GET("/userinfo", g.GetUserInfo) // 获取他人信息
		auth.POST("/report", g.Report)       // 举报
	}

	// 注册 handler
	service.Handle("/", router)

	err = service.Run()
	if err != nil {
		log.Error("服务启动失败: ", err)
		return
	}
}
