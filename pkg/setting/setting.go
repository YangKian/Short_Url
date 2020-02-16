package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type Server struct {
	RunMode      string
	HttpPort     int
	HeartBeatCheckTimer time.Duration
	ReadTimeoutTimer time.Duration
	WriteTimeoutTimer time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var DatabaseSetting = &Database{}

type IdGeneratorConfig struct {
	TimeStamp int64
	IdMaxBits int
	NodeMaxBits int
	Node int
}

var IdGeneratorSetting = &IdGeneratorConfig{}

var cfg *ini.File

func Start() {
	var err error
	if cfg, err = ini.Load("conf/app.ini"); err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v\n", err)
	}

	mapTo("database", DatabaseSetting)
	mapTo("server", ServerSetting)
	mapTo("idgenerator", IdGeneratorSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v) //将提取出来的配置映射到特定的数据结构中
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v\n", section, err)
	}
}
