package core

import (
	"fmt"
	"gitee.com/go-server/conf"
	"gitee.com/go-server/global"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const ConfigFile = "settings.yaml"

// InitConf 读取yaml文件配置
func InitConf() {
	c := &conf.Config{}
	yamlConf, err := os.ReadFile(ConfigFile)
	if err != nil {
		panic(fmt.Errorf("获取yaml配置文件错误:%s", err))
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("配置文件初始序列化失败：%v", err)
	}
	log.Println("配置文件加载成功.")
	// 把配置传送给全局的Config global文件中定义了Config
	global.Config = c
}
