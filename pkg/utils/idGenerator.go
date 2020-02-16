package utils

import (
	"errors"
	"sync"
	"time"
)

var (
	paramError = errors.New("init error, please check the params")
	timeError = errors.New("time setting error, Check your system time")
	)

type IdGenerator struct {
	lock     sync.Mutex
	lastTime int64 //上一秒生成 id 的时间
	count    int64 // 时间段内一共生成的 id 数
	node int // 当前节点的编号

	startTimeStamp int64 //基准时间戳
	maxIdBits int // 传入二进制位数，雪花算法的最后一段，配置1毫秒内能生成的最大 id 数，
	maxNodeBits int // 传入二进制位数，代表节点数

	MaxIdPerSecond int64
	MaxNode int64
	timeShift uint8
	nodeShift uint8
}

func NewIdGenerator(startTimeStamp int64, node, maxIdBits, maxNodeBits int) (*IdGenerator, error) {
	maxNode := -1 ^ (-1 << maxNodeBits)

	if err := checkParam(maxIdBits, maxNodeBits, node, maxNode); err != nil {
		return nil, err
	}

	return &IdGenerator{
		lock:     sync.Mutex{},
		lastTime: 0,
		count:    0,
		node: node,

		startTimeStamp: startTimeStamp,
		maxIdBits:      maxIdBits,
		maxNodeBits:    maxNodeBits,
		MaxIdPerSecond: -1 ^ (-1 << maxIdBits),
		MaxNode:        int64(maxNode),
		nodeShift:      uint8(maxIdBits),
		timeShift:      uint8(maxIdBits + maxNodeBits),
	}, nil
}

func (ig *IdGenerator) GetId() (int64, error) {
	ig.lock.Lock()
	defer ig.lock.Unlock()

	now := time.Now().UnixNano() / int64(time.Millisecond) //获取当前时间戳，转换为毫秒

	if now < ig.lastTime { // 检查是否出现时间戳回退的情况
		return -1, timeError
	}

	if now == ig.lastTime { // 判断当前时间是否还与上一次生成id的时间一致
		ig.count = (ig.count + 1) & ig.MaxIdPerSecond
		if ig.count == 0 { //检查当前毫秒内生成的id数是否已经超过限制
			for now <= ig.lastTime { //如果当前生成的 id 数已经超过上限，则需要等待到下一毫秒才能生成
				now = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else { //时间戳与上一次生成id的时间戳不同，说明已经不是同一毫秒了，则重置计数，并更新记录的时间戳
		ig.count = 0
	}

	ig.lastTime = now

	//最终id的构成：
	return (now - ig.startTimeStamp) << ig.timeShift | int64(ig.node) << ig.nodeShift |  ig.count, nil
}

func checkParam(maxIdBits, maxNodeBits, node, maxNode int) error {
	if maxIdBits < 0 || maxIdBits+maxNodeBits > 22 {
		return paramError
	}

	if maxNodeBits < 0 || maxNodeBits > 10 {
		return paramError
	}

	if node < 0 || node > maxNode {
		return paramError
	}
	return nil
}
