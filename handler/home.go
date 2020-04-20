package handler

import (
	"crypto/tls"
	"datalot/base"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Gin struct{}

const (
	name = "data_lot"
	desc = "世界很大我们曾相遇"
)

var (
	dataUrlVal = url.Values{}
	sid        int
)

func (g *Gin) Home(c *gin.Context) {

}

/**
	高德地图:
		1.首先需要创建一个轨迹服务
		2.在这个服务下添加终端，每个服务最多添加10万个
		3.查询轨迹
**/
func CreateService() {
	key := base.MapConf.Key
	data := make(map[string]string)
	data = map[string]string{
		"key":  key,
		"name": name,
		"desc": desc,
	}
	for key, val := range data {
		dataUrlVal.Add(key, val)
	}

	req, err := http.NewRequest("POST", base.MapConf.ServiceUrl, strings.NewReader(dataUrlVal.Encode()))
	if err != nil {
		log.Error("创建请求失败: ", err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return
	}

	res := make(map[string]interface{})
	_ = json.Unmarshal(body, &res)
	//fmt.Println("Result: ", res)
	for key, val := range res {
		if key == "data" {
			if val == nil{
				fmt.Println("Error: ",res)
				return
			}
			for k, v := range val.(map[string]interface{}) {
				if k == "sid" {
					sid = int(v.(float64))
				}
			}
		}
	}
	fmt.Println("Service ID: ", sid)
}

func addTerminal(meid, desc string) (int, error) {
	key := base.MapConf.Key
	var tid int
	serviceId := strconv.Itoa(sid)
	dataUrlVal = url.Values{}
	data := make(map[string]string)
	data = map[string]string{
		"key":  key,
		"sid":  serviceId,
		"name": meid,
		"desc": desc,
	}
	for key, val := range data {
		dataUrlVal.Add(key, val)
	}

	req, err := http.NewRequest("POST", base.MapConf.TerminalUrl, strings.NewReader(dataUrlVal.Encode()))
	if err != nil {
		log.Error("创建请求失败: ", err)
		return 0, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
		return 0, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	res := make(map[string]interface{})
	_ = json.Unmarshal(body, &res)
	for key, val := range res {
		if key == "data" {
			for k, v := range val.(map[string]interface{}) {
				if k == "tid" {
					tid = int(v.(float64))
				}
			}
		}
	}
	fmt.Println("Terminal ID: ", tid)
	return tid, nil
}

func addTrack(tid int) (int, error) {
	key := base.MapConf.Key
	var trid int
	serviceId := strconv.Itoa(sid)
	terminalId := strconv.Itoa(tid)
	dataUrlVal = url.Values{}
	data := make(map[string]string)
	data = map[string]string{
		"key": key,
		"sid": serviceId,
		"tid": terminalId,
	}
	for key, val := range data {
		dataUrlVal.Add(key, val)
	}
	req, err := http.NewRequest("POST", base.MapConf.TrackUrl, strings.NewReader(dataUrlVal.Encode()))
	if err != nil {
		log.Error("创建请求失败: ", err)
		return 0, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
		return 0, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	res := make(map[string]interface{})
	_ = json.Unmarshal(body, &res)
	for key, val := range res {
		if key == "data" {
			for k, v := range val.(map[string]interface{}) {
				if k == "trid" {
					trid = int(v.(float64))
				}
			}
		}
	}
	fmt.Println("Track ID: ", trid)
	return trid, nil
}

// 查询轨迹
