package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

var cfg *ini.File

func Setup() {
	var err error
	if Cfg, err = ini.Load("conf/app.ini"); err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v\n", err)
	}

	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)

	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v) //将提取出来的配置映射到特定的数据结构中
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v\n", section, err)
	}
}