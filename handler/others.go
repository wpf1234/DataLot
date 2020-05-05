package handler

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// 在线查看文档或视频
func (g *Gin) OnlinePreview(c *gin.Context) {
	var content string
	file := c.Query("file")

	// 读取文件
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("读取文件失败: ", err)
		c.JSON(200, gin.H{
			"code":    http.StatusInternalServerError,
			"data":    err,
			"message": "读取文件失败!",
		})
		return
	}

	switch {
	case strings.HasSuffix(file, ".pdf"):
		content = "application/pdf"
		break
	case strings.HasSuffix(file, ".mp4"):
		content = "video/mp4"
		break
	case strings.HasSuffix(file, ".avi"):
		content = "video/avi"
		break
	case strings.HasSuffix(file, ".jpeg") || strings.HasSuffix(file, ".jpg"):
		content = "image/jpeg"
		break
	case strings.HasSuffix(file, ".png"):
		content = "image/png"
		break
	case strings.HasSuffix(file, ".gif"):
		content = "image/gif"
		break
	default:
		content = "application/octet-stream"
	}
	c.Writer.Header().Add("Content-Type", content)
	c.Writer.Write(data)
}

