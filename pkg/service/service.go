package services

import (
	"MyProject/Short_Url/pkg/lru"

	"github.com/gin-gonic/gin"
)

type service struct {
	urlCache   *lru.LRUCache
	shortCache *lru.LRUCache
}

func Create(c *gin.Context) {
}
