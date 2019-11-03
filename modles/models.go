package models

import (
	"MyProject/Short_Url/pkg/setting"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func Start() {
	var err error

	db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models init fault: %v", err)
	}

	db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateForCreate)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateForUpdate)
	db.Callback().Delete().Replace("gorm:delete", updateForDelete)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

}

//创建时，同时添加创建时间和修改时间
func updateForCreate(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTime, ok := scope.FieldByName("CreatedOn"); ok {
			createTime.Set(nowTime)
		}

		if modifyTime, ok := scope.FieldByName("ModifiedOn"); ok {
			modifyTime.Set(nowTime)
		}
	}
}

//TODO:确认此处的逻辑
func updateForUpdate(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

//TODO:待完成
func updateForDelete(scope *gorm.Scope) {
	if !scope.HasError() {
	}
}
