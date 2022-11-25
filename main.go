package main

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

type TableAttr struct {
	ID                uuid.UUID `gorm:"column:id" sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Description       string    `gorm:"column:description" "form:"description" json:"description,omitempty"`
	Modifier          string    `gorm:"column:modifier" form:"modifier" json:"modifier,omitempty"`
	DeptBelongId      string    `gorm:"column:dept_belong_id" form:"dept_belong_id" json:"dept_belong_id,omitempty"`
	UpdateDateTime    time.Time `gorm:"column:update_datetime" form:"update_datetime" json:"update_datetime,omitempty"`
	CreateDateTime    time.Time `gorm:"column:create_datetime" form:"create_datetime" json:"create_datetime,omitempty"`
	TableUid          string    `gorm:"column:table_uid" form:"table_uid" json:"table_uid,omitempty"`
	AttrName          string    `gorm:"column:attr_name" form:"attr_name" json:"attr_name,omitempty"`
	Fixed             string    `gorm:"column:fixed" form:"fixed" json:"fixed,omitempty"`
	Prop              string    `gorm:"column:prop" form:"prop" json:"prop,omitempty"`
	Width             string    `gorm:"column:width" form:"width" json:"width,omitempty"`
	ShowSetting       int       `gorm:"column:showSetting" form:"showSetting" json:"showSetting,omitempty"`
	ColType           string    `gorm:"column:colType" form:"colType" json:"colType,omitempty"`
	ReturnNum         int       `gorm:"column:returnNum" form:"returnNum" json:"returnNum,omitempty"`
	LocalModel        string    `gorm:"column:localModel" form:"localModel" json:"localModel,omitempty"`
	NetModel          string    `gorm:"column:netModel" form:"netModel" json:"netModel,omitempty"`
	ObjectName        string    `gorm:"column:object_name" form:"object_name" json:"object_name,omitempty"`
	FieldName         string    `gorm:"column:field_name" form:"field_name" json:"field_name,omitempty"`
	RegularExpression string    `gorm:"column:regular_expression" form:"regular_expression" json:"regular_expression,omitempty"`
	FilterThredshold  float32   `gorm:"column:filter_threshold" form:"filter_threshold" json:"filter_threshold,omitempty"`
	ShowIndex         int       `gorm:"column:show_index" form:"show_index" json:"show_index,omitempty"`
	IncludeEntityType bool      `gorm:"column:include_entity_type" form:"include_entity_type" json:"include_entity_type,omitempty"`
}

func (v TableAttr) TableName() string {
	return "fill_table_attr_info"
}

func main() {
	db, err := gorm.Open("mysql", "root:Lettcue2kg@tcp(192.168.1.247:28898)/atomecho2?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("连接数据库失败")
	}
	logrus.Info("数据库连接成功")

	defer db.Close()
	// 查找
	table := TableAttr{}
	db.First(&table, "table_uid = ?", "191cbac6-cb89-311a-a5e0-2f92e52aa4ef")
	logrus.Info("table:", table.AttrName)

	// // 创建
	// new_info := TableAttr{
	// 	ID: uuid.NewV5(uuid.NamespaceURL, "1ddasdasd2"),
	// }
	// db.Create(&new_info) // 通过数据的指针来创建

	// // 更新
	// db.First(&table, "table_uid = ?", "191cbac6-cb89-311a-a5e0-2f92e52aa4ef")
	// table.AttrName = "111"
	// db.Model(&table).Updates(table)

	// // 删除 - 删除product
	// db.Delete(&table)
}
