package models

import (
	"MyProject/Short_Url/pkg/utils"
	"fmt"

	"github.com/jinzhu/gorm"
)

type UrlCode struct {
	Model

	MD5    string
	Code   string
	Url    string
	Click  int
	UserId int
}

//添加一个url
func (u UrlCode) AddUrl(url string, userId int) (int, error) {
	urlCode := UrlCode{
		MD5:    utils.MD5(url),
		Code:   "",
		UserID: userId,
	}

	if err := db.Create(&urlCode).Error; err != nil {
		fmt.Printf("url添加失败, err: %v\n", err)
		return 0, err
	}

	return urlCode.ID, nil
}

//TODO: 返回值设什么好？
func (u UrlCode) GetByCode(code string) UrlCode {
	var res UrlCode
	if err := db.Where("code = ?", code).First(&res).Error; err != nil {
		fmt.Printf("GetByCode failed, err: %v\n", err)
		return nil
	}

	return res
}

func (u UrlCode) UpdateCode(id int, code string) error {
	db.Table("url_codes").Where("id = ?", id).Update("code", code)
	if db.Error != nil {
		fmt.Printf("url update faile, err: %v\n", db.Error)
		return db.Error
	}
	return nil
}

//TODO:写log？
func (u UrlCode) SaveClicks(clicks map[string]int) {
	for code, c := range clicks {
		db.Where("code = ?", code).Find(&UrlCode{}).UpdateColumn("click", gorm.Expr("click + ?", c))
		fmt.Println("add %d click on %s", c, code)
	}
}
