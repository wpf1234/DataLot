package base

import (
	"datalot/models"
	"datalot/utils"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	log "github.com/sirupsen/logrus"
)

const defaultPath = "app"

var (
	m   sync.RWMutex
	key = []byte("data_lot*2020key")

	mc      models.MysqlConf
	lc      models.LogConf
	//MapConf models.MapConf
	IMC models.IMConf
	MapBox models.MapBox

	DB *gorm.DB
)

func Init() {
	m.Lock()
	defer m.Unlock()

	err := config.Load(file.NewSource(
		file.WithPath("./conf/application.yml"),
	))

	if err != nil {
		log.Error("加载配置文件失败: ", err)
		return
	}

	if err := config.Get(defaultPath, "log").Scan(&lc); err != nil {
		log.Error("Log 配置读取失败: ", err)
		return
	}
	log.Info("Log 配置读取成功!")
	utils.LoggerToFile(lc.LogPath, lc.LogFile)

	if err := config.Get(defaultPath, "mysql").Scan(&mc); err != nil {
		log.Error("Mysql 配置文件读取失败: ", err)
		return
	}
	log.Info("读取 Mysql 配置成功!")

	str, err := base64.StdEncoding.DecodeString(mc.Password)
	if err != nil {
		log.Error("Base64 decode failed: ", err)
		return
	}
	pwd, err := utils.AesDecrypt(str, key)
	if err != nil {
		log.Error("ASE decrypt failed: ", err)
		return
	}
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=true",
		mc.User, string(pwd), mc.Host, mc.DB)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Error("Open mysql failed: ", err)
		return
	}

	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(10)
	DB.DB().SetConnMaxLifetime(10 * time.Second)

	if err := DB.DB().Ping(); err != nil {
		log.Error("连接数据库失败: ", err)
		return
	}
	log.Info("数据库连接成功!")

	if err := config.Get(defaultPath, "map").Scan(&MapBox); err != nil {
		log.Error("地图资源配置文件读取失败: ", err)
		return
	}
	log.Info("读取地图资源配置成功!")

	if err := config.Get(defaultPath, "tencent").Scan(&IMC); err != nil {
		log.Error("IM即时通讯配置失败: ", err)
		return
	}
	log.Info("IM即时通讯配置成功!")
}
