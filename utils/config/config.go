package config

import (
	"douyu/utils/db"
	"douyu/utils/redis"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 加载配置文件
func InitConfig(configFile *string) {
	splits := strings.Split(filepath.Base(*configFile), ".")
	viper.SetConfigName(filepath.Base(splits[0]))
	viper.AddConfigPath(filepath.Dir(*configFile))
	viper.SetConfigType(filepath.Ext(splits[0]))
	if err := viper.ReadInConfig(); err != nil {
		log.Println("加载配置文件失败;异常信息:"+err.Error())
		time.Sleep(time.Second * 3)
		os.Exit(0)
	}
	// 使用配置
	if err := ConnResources(); err != nil {
		log.Println("应用配置文件失败;异常信息:"+err.Error())
		time.Sleep(time.Second * 3)
		os.Exit(0)
	}
	// 配置监听
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("配置文件变更")
		if err := ConnResources(); err != nil {
			log.Println("重新加载配置文件失败;异常信息:"+err.Error())
		}
	})
	log.Println( "成功加载配置文件,并监听变更中...")
}

func ConnResources() error {
	// 重载配置
	str := "启动失败!"
	err := db.InitWithEnv()
	if err != nil {
		return errors.New(str + "数据库连接异常:" + err.Error())
	}
	err = redis.InitWithEnv()
	if err != nil {
		return errors.New(str + "缓存连接异常:" + err.Error())
	}
	return nil
}
