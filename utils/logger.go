package utils

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func LoggerToFile(logPath, logFile string) {
	// 日志文件
	fileName := path.Join(logPath, logFile)
	// 写入文件
	//src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	//if err != nil {
	//	fmt.Println("日志写入失败: ", err)
	//	return
	//}
	// 设置输出
	log.SetOutput(os.Stdout)
	// 设置日志级别
	log.SetLevel(log.DebugLevel)
	// 设置打印该条日志的函数名及行号
	log.SetReportCaller(true)
	// 设置 rotate log
	logWrite, err := rotatelogs.New(
		// 分割后的文件名
		fileName+".%Y%m%d.log",
		//生成软链
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志分割时间间隔
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		fmt.Println("Rotate log error: ", err)
		return
	}
	writeMap := lfshook.WriterMap{
		log.InfoLevel:  logWrite,
		log.FatalLevel: logWrite,
		log.DebugLevel: logWrite,
		log.WarnLevel:  logWrite,
		log.ErrorLevel: logWrite,
		log.PanicLevel: logWrite,
	}

	lfHook := lfshook.NewHook(writeMap, &log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.AddHook(lfHook)
}

