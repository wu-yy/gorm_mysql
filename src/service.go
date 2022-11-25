package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type TotalEnergy struct {
	Id          int       `gorm:"column:id" form:"id" json:"id,omitempty"`
	GmtCreate   time.Time `gorm:"column:gmt_create" form:"gmt_create" json:"gmt_create,omitempty"`
	GmtModified time.Time `gorm:"column:gmt_modified" form:"gmt_modified" json:"gmt_modified,omitempty"`
	UserId      string    `gorm:"column:user_id" form:"user_id" json:"user_id,omitempty"`
	TotalEnergy int       `gorm:"column:total_energy" form:"total_energy" json:"total_energy,omitempty"`
}

func (v TotalEnergy) TableName() string {
	return "total_energy"
}

type ToCollectEnergy struct {
	Id              int       `gorm:"column:id" form:"id" json:"id,omitempty"`
	GmtCreate       time.Time `gorm:"column:gmt_create" form:"gmt_create" json:"gmt_create,omitempty"`
	GmtModified     time.Time `gorm:"column:gmt_modified" form:"gmt_modified" json:"gmt_modified,omitempty"`
	UserId          string    `gorm:"column:user_id" form:"user_id" json:"user_id,omitempty"`
	ToCollectEnergy int       `gorm:"column:to_collect_energy" form:"to_collect_energy" json:"to_collect_energy,omitempty"`
	Status          string    `gorm:"column:status" form:"status" json:"status,omitempty"`
}

func (v ToCollectEnergy) TableName() string {
	return "to_collect_energy"
}

type Hub struct {
	SqlDb *gorm.DB
}

var single_hub *Hub

func NewHub() *Hub {
	logrus.Info("new hub")
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",
		"111111",
		"127.0.0.1",
		3306,
		"atec2022"))
	if err != nil {
		logrus.Error("connect to database error: ", err)
	}
	single_hub = &Hub{
		SqlDb: db,
	}
	return single_hub
}

func GetSinletonHub() *Hub {
	if single_hub != nil {
		return single_hub
	}
	logrus.Fatal("Hub is not available")
	return nil
}

// 判断是不是和i属于用户
func IsBelongUser(user_id string, collect_id int) bool {
	table := ToCollectEnergy{}
	GetSinletonHub().SqlDb.First(&table, "id = ?", collect_id)
	// 属于用户
	if table.UserId == user_id {
		// 更新to_collect_energy
		enery_num := table.ToCollectEnergy
		table.Status = "all_collected"
		GetSinletonHub().SqlDb.Update(table)

		// 更新用户的能量
		user_energy := TotalEnergy{}
		GetSinletonHub().SqlDb.First(&user_energy, "user_id = ?", user_id)
		user_energy.TotalEnergy += enery_num
		GetSinletonHub().SqlDb.Update(user_energy)
		return true
	} else { // 不属于用户
		// 不属于用户被采集过
		if table.Status == "collected_by_other" || table.Status == "all_collected" {
			return true
		}
		// 更新to_collect_energy
		enery_num := table.ToCollectEnergy
		table.Status = "collected_by_other"
		table.ToCollectEnergy -= int(math.Floor(float64(enery_num) * 0.3))
		GetSinletonHub().SqlDb.Update(table)

		// 更新用户的能量
		user_energy := TotalEnergy{}
		GetSinletonHub().SqlDb.First(&user_energy, "user_id = ?", user_id)
		user_energy.TotalEnergy += int(math.Floor(float64(enery_num) * 0.3))
		GetSinletonHub().SqlDb.Update(user_energy)
		return true
	}
}

func collect_energy(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if len(ps) == 2 {
		logrus.Info("userid: %s collectid:%s", ps[0].Value, ps[1].Value)
		var user_id string = ps[0].Value
		var collected_id int
		num, err := strconv.ParseInt(ps[1].Value, 10, 32)
		if err != nil {
			fmt.Fprintf(w, "false")
		}
		collected_id = int(num)
		// 逻辑函数
		IsBelongUser(user_id, collected_id)
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}

}

func main() {
	NewHub()
	router := httprouter.New()
	router.GET("/collect_energy/:userid/:collectid", collect_energy)
	http.ListenAndServe(":18090", router)
}
