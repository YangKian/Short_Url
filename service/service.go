package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"shortUrl/constants"
	"shortUrl/models"
	"shortUrl/pkg/lru"
	"shortUrl/pkg/setting"
	"shortUrl/pkg/utils"
	"sync"
)

var clickCounter counter

type counter struct {
	counterMap map[string]int
	flag       bool
	lock       sync.Mutex
}

func newCounter() counter {
	return counter{
		counterMap: make(map[string]int),
		flag:       false,
	}
}

func Start() {
	clickCounter = newCounter()
}

type Service interface {
	SingleCreate(*gin.Context)
	TransToUrl(*gin.Context)
}

type shortUrlService struct {
	urlCache    *lru.LRUCache
	shortCache  *lru.LRUCache
	lock        sync.RWMutex
	idGenerator *utils.IdGenerator
}

func NewShortUrlService() Service {
	idGenerator, err := utils.NewIdGenerator(setting.IdGeneratorSetting.TimeStamp,
		setting.IdGeneratorSetting.Node,
		setting.IdGeneratorSetting.IdMaxBits,
		setting.IdGeneratorSetting.NodeMaxBits)
	if err != nil {
		log.Fatalf("start service err: %v\n", err)
	}
	return &shortUrlService{
		urlCache:    lru.Constructors(100),
		shortCache:  lru.Constructors(100),
		idGenerator: idGenerator,
	}
}

func (s *shortUrlService) SingleCreate(c *gin.Context) {
	url := c.PostForm("url")
	fmt.Printf("[Create]: get url from gin context: %s\n", url)
	fmt.Printf("[Create]: Begin to create a short url, the given url is: %s\n", url)

	s.getShortCode(c, url)
}

func (s *shortUrlService) getShortCode(c *gin.Context, url string) {
	s.lock.RLock()
	res, err := s.urlCache.Get(url, constants.ORIGINURL)
	s.lock.RUnlock()

	if err != nil {
		//数据库和缓存中都没有该url的记录，则新建shortCode
		switch err {
		case gorm.ErrRecordNotFound:
			fmt.Printf("[Create]: 缓存和数据库中未查到url，新建shortCode, url: %s\n", url)
			//生成short url并更新数据库
			shortCode, err := s.codeGenerator(url)
			if err != nil {
				code := constants.CREATE_shortUrl_ERROR
				innerFail(c, code)
				log.Println("[Create]: 创建短地址失败")
				return
			}

			go func() {
				s.lock.Lock()
				s.urlCache.Put(url, shortCode)
				s.shortCache.Put(shortCode, url)
				s.lock.Unlock()
			}()

			code := constants.SUCCESS
			success(c, code, shortCode)
			log.Printf("[Create]: create shorturl for url success, shorturl is: %s\n", shortCode)
			return
		default:
			log.Printf("[getShortCode]：运行出错，err: %v\n", err)
			code := constants.DBERROER
			innerFail(c, code)
			return
		}
	}

	code := constants.SUCCESS
	success(c, code, res)
	log.Printf("[Create]: get shorturl for url success, shorturl is: %s\n", res)
	return
}

func (s *shortUrlService) TransToUrl(c *gin.Context) {
	fmt.Println(c.Request)
	shortUrl := c.PostForm("shortUrl")
	fmt.Printf("[TransToUrl]: get shortUrl from gin context: %s\n", shortUrl)

	s.lock.RLock()
	url, err := s.shortCache.Get(shortUrl, constants.SHORTURL)
	s.lock.RUnlock()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			code := constants.NOTFOUND
			requestFail(c, code)
			return
		}

		code := constants.REQUEST_ERROR
		requestFail(c, code)
		return
	}

	visitCount(shortUrl)
	redirect(c, url)
	log.Printf("[TransToUrl]: get url success, url: %s\n", url)
	return
}

func (s *shortUrlService) codeGenerator(url string) (string, error) {
	urlCode := models.UrlCode{}
	urlId, err := s.idGenerator.GetId() //为当前url设置一个唯一id
	if err != nil {
		return "", err
	}
	//生成短域名
	shortCode := utils.Transport(urlId)
	//更新数据库
	if err := urlCode.AddUrl(url, shortCode); err != nil {
		return "", err
	}
	return shortCode, nil
}

func visitCount(shortUrl string) {
	clickCounter.lock.Lock()
	defer clickCounter.lock.Unlock()

	if _, ok := clickCounter.counterMap[shortUrl]; ok {
		clickCounter.counterMap[shortUrl]++
	} else {
		clickCounter.counterMap[shortUrl] = 1
	}
	clickCounter.flag = true
}
