package utils

import (
	"sync"
	"time"
)

const (
	maxIdPerSecond       = -1 ^ (-1 << 12) // 对应雪花算法的最后一段，12个bit，说明1毫秒内能生成的最大id数
	startTimeStamp       = 1576818924228   // 起始时间戳
	timeShift      uint8 = 12 + 10
)

type IdGenerator struct {
	lock     sync.Mutex
	lastTime int64
	count    int
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		lock:     sync.Mutex{},
		lastTime: 0,
		count:    0,
	}
}

func (ig *IdGenerator) GetId() int64 {
	ig.lock.Lock()
	defer ig.lock.Unlock()

	now := time.Now().UnixNano() / int64(time.Millisecond) //获取当前时间戳，转换为毫秒
	if now == ig.lastTime {                                //判断当前时间是否还与上一次生成id的时间一致
		ig.count++
		if ig.count > maxIdPerSecond { //检查当前毫秒内生成的id数是否已经超过限制
			time.Sleep(time.Millisecond) //如果超过，则等待一毫秒，之后再继续生成
			now = time.Now().UnixNano() / int64(time.Millisecond)
		}
	} else { //时间戳与上一次生成id的时间戳不同，说明已经不是同一毫秒了，则重置计数，并更新记录的时间戳
		ig.count = 0
		ig.lastTime = now
	}

	//最终id的构成：
	return (now - startTimeStamp) << timeShift | int64(1 << 12) | int64(ig.count)
}
