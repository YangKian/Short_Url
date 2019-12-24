package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

type UrlCode struct {
	Model

	Code  string `gorm:"code"`
	Url   string `gorm:"url"`
	Click int    `gorm:"clikc"`
}

//添加一个url
func (u UrlCode) AddUrl(url string, code string) error {
	urlCode := UrlCode{
		Url:  url,
		Code: code,
	}

	if err := db.Create(&urlCode).Error; err != nil {
		log.Printf("url添加失败, err: %v\n", err)
		return err
	}

	return nil
}

//根据code查询
func (u UrlCode) GetByCode(code string) (UrlCode, error) {
	var res UrlCode
	if err := db.Table("url_code").Where("code = ?", code).First(&res).Error; err != nil {
		log.Printf("GetByCode failed, code: %v, err: %v\n", code, err)
		return UrlCode{}, err
	}

	return res, nil
}

//根据url查询
func (u UrlCode) GetByUrl(url string) (UrlCode, error) {
	var res UrlCode
	if err := db.Table("url_code").Where("url = ?", url).First(&res).Error; err != nil {
		log.Printf("GetByUrl failed, err: %v\n", err)
		return UrlCode{}, err
	}

	return res, nil
}

func (u UrlCode) UpdateCode(id int, code string) error {
	fmt.Printf("[UpdateCode]: update code, id: %d, code: %s\n", id, code)
	tx := db.Begin()
	tx.Table("url_code").Where("id = ?", id).Update("code", code)
	if db.Error != nil {
		tx.Rollback()
		log.Printf("[UpdateCode]: url update faile, err: %v\n", db.Error)
		return db.Error
	}
	tx.Commit()
	log.Printf("[UpdateCode]: url update success, code: %v\n", code)
	return nil
}

func (u UrlCode) SaveClicks(clicks map[string]int) {
	for code, c := range clicks {
		db.Table("url_code").Where("code = ?", code).Find(&UrlCode{}).
			UpdateColumn("click", gorm.Expr("click + ?", c))
		log.Printf("add %d click on %s\n", c, code)
	}
}
