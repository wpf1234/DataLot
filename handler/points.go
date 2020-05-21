package handler

import (
	"datalot/base"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func InitPoints() {
	db := base.DB.Raw("select id from user")
	rows, err := db.Rows()
	if err != nil {
		log.Error("查询失败: ", err)
		return
	}

	for rows.Next() {
		var id int
		_ = rows.Scan(&id)

		now := time.Now()
		begin := now.UnixNano() / 10e5
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		end := next.UnixNano() / 10e5
		var pId int
		db=base.DB.Raw("select id from point where u_id=? and end=?",id,end)
		_=db.Row().Scan(&pId)
		if pId == 0{
			db = base.DB.Exec("insert into point set u_id=?,begin=?,end=?", id, begin, end)
			fmt.Println("Insert: ", db.RowsAffected)
		}

	}
	_ = rows.Close()
}
