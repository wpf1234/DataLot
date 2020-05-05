package handler

import (
	"datalot/base"
	"datalot/models"
	"datalot/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (g *Gin) HeadPortrait(c *gin.Context) {
	token := c.Request.Header.Get("token")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("获取请求失败!")

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    err,
			"message": "获取数据失败!",
		})
		return
	}

	// 获取文件名
	fileName := header.Filename
	str := strings.Split(fileName, ".")
	layout := strings.ToLower(str[len(str)-1])
	if layout != "jpeg" && layout != "png" && layout != "jpg" && layout != "gif" {
		log.Error("文件格式不正确!")
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    nil,
			"message": "文件格式不正确!",
		})
		return
	}

	if header.Size > 2000000 {
		//判断大小是否大于2M
		log.Error("文件过大!")
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    nil,
			"message": "文件大于2M，请重新上传",
		})
		return
	}
	claims, err := utils.ParseToken(token)
	if err != nil {
		log.Error("Token解析失败: ", err)
		return
	}
	id := claims.Id
	username := claims.Username
	filePath := "static/head/" + username + "-" + fileName
	exist, err := utils.PathExists(filePath)
	if err != nil {
		log.Error(err)
		return
	}
	// 存在
	if exist {
		return
	}
	// 不存在
	out, err := os.Create(filePath)
	if err != nil {
		log.Error("创建文件失败!")
		c.JSON(200, gin.H{
			"code":    500,
			"data":    err,
			"message": "创建文件失败!",
		})
		return
	}
	_, err = io.Copy(out, file)
	if err != nil {
		log.Error(err)
		c.JSON(200, gin.H{
			"code":    500,
			"data":    err,
			"message": "保存文件失败!",
		})
		return
	}
	file.Close()
	out.Close()

	db := base.DB.Exec("update user set head_portrait=? where id=?", filePath, id)
	log.Info("更新: ", db.RowsAffected)

	c.JSON(200, gin.H{
		"code":    200,
		"data":    filePath,
		"message": "上传图片成功!",
	})

}

func (g *Gin) CyclePicture(c *gin.Context) {
	var res models.UploadPIC
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

	tm := strconv.FormatInt(time.Now().Unix(), 10)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("获取请求失败!")

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    err,
			"message": "获取数据失败!",
		})
		return
	}

	// 获取文件名
	fileName := header.Filename
	str := strings.Split(fileName, ".")
	layout := strings.ToLower(str[len(str)-1])
	if layout != "jpeg" && layout != "png" && layout != "jpg" && layout != "gif" {
		log.Error("文件格式不正确!")
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    nil,
			"message": "文件格式不正确!",
		})
		return
	}

	if header.Size > 10000000 {
		//判断大小是否大于10M
		log.Error("文件过大!")
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"data":    nil,
			"message": "文件大于2M，请重新上传",
		})
		return
	}
	dir := "static/dynamic/" + username
	//dir := "static/dynamic/test"
	dirExit, err := utils.PathExists(dir)
	if err != nil {
		log.Error("Dir error: ", err)
		return
	}
	if !dirExit {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Error("创建目录失败: ", err)
			c.JSON(http.StatusOK, gin.H{
				"code":    http.StatusInternalServerError,
				"data":    nil,
				"message": "创建用户目录失败!",
			})
			return
		}
	}
	filePath := dir + "/" + tm + "_" + fileName
	out, err := os.Create(filePath)
	if err != nil {
		log.Error("创建文件失败!")
		c.JSON(200, gin.H{
			"code":    500,
			"data":    err,
			"message": "创建文件失败!",
		})
		return
	}
	_, err = io.Copy(out, file)
	if err != nil {
		log.Error(err)
		c.JSON(200, gin.H{
			"code":    500,
			"data":    err,
			"message": "保存文件失败!",
		})
		return
	}
	file.Close()
	out.Close()

	//fhs := c.Request.MultipartForm.File["image"]
	//fmt.Println("12312312312312: ",fhs)
	//for _, fh := range fhs {
	//	if fh.Size > 20000000 {
	//		log.Warn("图片大于20M!!!")
	//		c.JSON(http.StatusOK, gin.H{
	//			"code":    http.StatusBadRequest,
	//			"data":    nil,
	//			"message": "文件大于20M，请重新上传",
	//		})
	//		return
	//	}
	//	file, _ := fh.Open()
	//	fileName := fh.Filename
	//	dir := "static/dynamic/" + username
	//	dirExit, err := utils.PathExists(dir)
	//	if err != nil {
	//		log.Error("Dir error: ", err)
	//		return
	//	}
	//	if !dirExit {
	//		err = os.Mkdir(dir, os.ModePerm)
	//		if err != nil {
	//			log.Error("创建目录失败: ", err)
	//			c.JSON(http.StatusOK, gin.H{
	//				"code":    http.StatusInternalServerError,
	//				"data":    nil,
	//				"message": "创建用户目录失败!",
	//			})
	//			return
	//		}
	//	}
	//
	//	filePath := dir + "/" + fileName + "_" + tm
	//	out, err := os.Create(filePath)
	//	if err != nil {
	//		log.Error("创建文件失败!")
	//		c.JSON(200, gin.H{
	//			"code":    500,
	//			"data":    err,
	//			"message": "创建文件失败!",
	//		})
	//		return
	//	}
	//	_, err = io.Copy(out, file)
	//	if err != nil {
	//		log.Error(err)
	//		c.JSON(200, gin.H{
	//			"code":    500,
	//			"data":    err,
	//			"message": "保存文件失败!",
	//		})
	//		return
	//	}
	//	file.Close()
	//	out.Close()
	//
	//	files = append(files, filePath)
	//}

	res.UserId = userId
	res.Picture = filePath

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    res,
		"message": "上传图片成功!",
	})
}
