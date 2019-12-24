package models

import (
	"fmt"
	"log"
	"shortUrl/pkg/setting"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on" gorm:"created_on"`
	ModifiedOn int `json:"modified_on" gorm:"modified_on"`
}

func Start() {
	var err error
	db, err = gorm.Open(setting.DatabaseSetting.Type,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Password,
			setting.DatabaseSetting.Host,
			setting.DatabaseSetting.Name))

	if err != nil {
		log.Fatalf("models init fault: %v", err)
	}

	db.SingularTable(true)
	db.LogMode(false)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateForCreate)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateForUpdate)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

}

//创建时，同时添加创建时间和修改时间
func updateForCreate(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTime, ok := scope.FieldByName("CreatedOn"); ok {
			err := createTime.Set(nowTime)
			if err != nil {
				fmt.Printf("[updateForCreate]: createTime err: %s\n", err)
			}
		}

		if modifyTime, ok := scope.FieldByName("ModifiedOn"); ok {
			err := modifyTime.Set(nowTime)
			if err != nil {
				fmt.Printf("[updateForCreate]: modifyTime err: %s\n", err)
			}
		}
	}
}

func updateForUpdate(scope *gorm.Scope) {
	if modifyTime, ok := scope.FieldByName("ModifiedOn"); ok {
		err := modifyTime.Set(time.Now().Unix())
		if err != nil {
			fmt.Printf("[updateForUpdate]: modifyTime err: %s\n", err)
		} else {
			fmt.Printf("[updateForUpdate]: modifyTime success")
		}
	}
}
