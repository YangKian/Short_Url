package service

import (
	"MyProject/Short_Url/contants"
	"MyProject/Short_Url/models"
	"MyProject/Short_Url/pkg/lru"
	"MyProject/Short_Url/pkg/utils"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	ORIGINURL = 1
	SHORTURL = 2
)

type Service interface {
	SingleCreate(*gin.Context)
	MultiCreate(*gin.Context)
}

type shortUrlService struct {
	urlCache   *lru.LRUCache
	shortCache *lru.LRUCache
	lock       sync.RWMutex
}

func NewShortUrlService() Service {
	return &shortUrlService{
		urlCache:   lru.Constructors(100),
		shortCache: lru.Constructors(100),
	}
}

func (s *shortUrlService) SingleCreate(c *gin.Context) {
	url := c.PostForm("url")
	fmt.Printf("[Create]: get url from gin context: %s\n", url)
	fmt.Printf("[Create]: Begin to create a short url, the given url is: %s\n", url)

	s.getShortCode(c, url)

	//s.lock.RLock()
	//res, err := s.urlCache.Get(url)
	//s.lock.RUnlock()
	//
	//if err != nil {
	//	//数据库和缓存中都没有该url的记录，则新建shortCode
	//	fmt.Printf("[Create]: 缓存和数据库中未查到url，新建shortCode, url: %s\n", url)
	//	if err == gorm.ErrRecordNotFound {
	//		userId := c.GetInt("userId")
	//		//生成short url并更新数据库
	//		shortCode, err := codeGenerator(url, userId)
	//		if err != nil {
	//			code := contants.CREATE_SHORT_URL_ERROR
	//			innerFail(c, code)
	//			fmt.Println("[Create]: 创建短地址失败")
	//			return
	//		}
	//
	//		go func() {
	//			s.lock.Lock()
	//			s.urlCache.Put(url, shortCode)
	//			s.shortCache.Put(shortCode, url)
	//			s.lock.Unlock()
	//		}()
	//
	//		code := contants.SUCCESS
	//		success(c, code, shortCode)
	//		fmt.Printf("[Create]: create shorturl for url success, shorturl is: %s\n", shortCode)
	//		return
	//	}
	//
	//	code := contants.DBERROER
	//	innerFail(c, code)
	//	return
	//}
	//
	//code := contants.SUCCESS
	//success(c, code, res)
	//fmt.Printf("[Create]: get shorturl for url success, shorturl is: %s\n", res)
	//return
}

func (s *shortUrlService) MultiCreate(c *gin.Context) {
	var body []string
	if err := c.ShouldBindJSON(&body); err != nil {
		code := contants.REQUEST_ERROR
		requestFail(c, code)
		fmt.Printf("请求错误， err: %s\n", err)
		return
	}

	if len(body) == 0 {
		code := contants.EMPTYREQUESTBODY
		requestFail(c, code)
		fmt.Println("请求体为空")
		return
	}

	for _, url := range body {
		go func(u string) {
			s.getShortCode(c, u)
		}(url)
	}

}

func (s *shortUrlService) getShortCode(c *gin.Context, url string) {
	s.lock.RLock()
	res, err := s.urlCache.Get(url)
	s.lock.RUnlock()

	if err != nil {
		//数据库和缓存中都没有该url的记录，则新建shortCode
		fmt.Printf("[Create]: 缓存和数据库中未查到url，新建shortCode, url: %s\n", url)
		if err == gorm.ErrRecordNotFound {
			userId := c.GetInt("userId")
			//生成short url并更新数据库
			shortCode, err := codeGenerator(url, userId)
			if err != nil {
				code := contants.CREATE_SHORT_URL_ERROR
				innerFail(c, code)
				fmt.Println("[Create]: 创建短地址失败")
				return
			}

			go func() {
				s.lock.Lock()
				s.urlCache.Put(url, shortCode)
				s.shortCache.Put(shortCode, url)
				s.lock.Unlock()
			}()

			code := contants.SUCCESS
			success(c, code, shortCode)
			fmt.Printf("[Create]: create shorturl for url success, shorturl is: %s\n", shortCode)
			return
		}

		code := contants.DBERROER
		innerFail(c, code)
		return
	}

	code := contants.SUCCESS
	success(c, code, res)
	fmt.Printf("[Create]: get shorturl for url success, shorturl is: %s\n", res)
	return
}

func codeGenerator(url string, userId int) (string, error) {
	urlCode := models.UrlCode{}
	urlId, err := urlCode.AddUrl(url, userId)
	if err != nil {
		return "", err
	}

	shortCode := utils.Transport(urlId)

	err = urlCode.UpdateCode(urlId, shortCode)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

func (s *shortUrlService) TransToUrl(c *gin.Context) {
	shortUrl := c.PostForm("shortUrl")

	fmt.Printf("[TransToUrl]: get shortUrl from gin context: %s\n", shortUrl)
	url, err := s.shortCache.Get(shortUrl)
	if err != nil {

	}
}


