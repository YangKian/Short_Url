package models

import (
	"MyProject/Short_Url/pkg/utils"
	"fmt"

	"github.com/jinzhu/gorm"
)

type UrlCode struct {
	Model

	MD5    string `gorm:"md5"`
	Code   string `gorm:"code"`
	Url    string `gorm:"url"`
	Click  int `gorm:"clikc"`
	UserId int `gorm:"user_id"`
}

//添加一个url
func (u UrlCode) AddUrl(url string, userId int) (int, error) {
	urlCode := UrlCode{
		Url:    url,
		MD5:    utils.MD5(url),
		Code:   "",
		UserId: userId,
	}

	if err := db.Create(&urlCode).Error; err != nil {
		fmt.Printf("url添加失败, err: %v\n", err)
		return 0, err
	}

	return urlCode.ID, nil
}

//根据code查询
func (u UrlCode) GetByCode(code string) (UrlCode, error) {
	var res UrlCode
	if err := db.Where("code = ?", code).First(&res).Error; err != nil {
		fmt.Printf("GetByCode failed, err: %v\n", err)
		return UrlCode{}, err
	}

	return res, nil
}

//根据url查询
func (u UrlCode) GetByUrl(url string) (UrlCode, error) {
	var res UrlCode
	if err := db.Where("url = ?", url).First(&res).Error; err != nil {
		fmt.Printf("GetByUrl failed, err: %v\n", err)
		return UrlCode{}, err
	}

	return res, nil
}

func (u UrlCode) UpdateCode(id int, code string) error {
	fmt.Printf("[UpdateCode]: update code, id: %s, code: %s\n", id, code)
	db.Table("url_codes").Where("id = ?", id).Update("code", code)
	if db.Error != nil {
		fmt.Printf("url update faile, err: %v\n", db.Error)
		return db.Error
	}
	fmt.Printf("url update success, code: %v\n", code)
	return nil
}

//TODO:写log？
func (u UrlCode) SaveClicks(clicks map[string]int) {
	for code, c := range clicks {
		db.Where("code = ?", code).Find(&UrlCode{}).UpdateColumn("click", gorm.Expr("click + ?", c))
		fmt.Println("add %d click on %s", c, code)
	}
}
