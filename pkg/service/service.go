package services

import (
	"MyProject/Short_Url/pkg/lru"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Sercice interface {
	Create(*gin.Context)
}

type shortUrlService struct {
	urlCache   *lru.LRUCache
	shortCache *lru.LRUCache
}

func NewShortUrlService() Service {
	return &shortUrlService{
		urlCache:   lru.Constructors(100),
		shortCache: lru.Constructors(100),
	}
}

func (s *shortUrlService) Create(c *gin.Context) {
	c.PostForm("url")
	fmt.Printf("Begin to create a short url, the given url is: %s\n", url)
	res, err := s.urlCache.Get(url)
	if err != nil {
		//数据库和缓存中都没有该url的记录，则新建shortCode
		if err == ErrRecordNotFound {
			md5 := utils.MD5(url)
			userId := c.GetInt("userId")
			shortCode, err := utils.CodeGenerator(url, userId)
			if err != nil {
				fmt.Println("创建短地址失败")
				c.JSON(http.StatusOK, gin.H{
					"code": 302,
					"msg":  "创建短地址失败",
				})
				return
			}

			fmt.Printf("create shorturl for url success, shorturl is: %s\n", shortCode)
			// s.urlCache.Put(url,
		}
	}
}

// }
