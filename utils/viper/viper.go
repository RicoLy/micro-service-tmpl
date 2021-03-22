package viper

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	ViperConf *viper.Viper
)

func init() {
	confPath := "./etc/"
	InitConfig(confPath)
}

// 初始化
func InitConfig(confPath string) {
	// 读取toml配置文件 设置配置文件名为 config, 不需要配置文件扩展名，配置文件的类型 viper 会自动根据扩展名自动匹配.
	ViperConf = viper.New()

	ViperConf.SetConfigName("Ai")      // 配置文件名称
	ViperConf.AddConfigPath(confPath)  //设置配置文件的搜索目录
	ViperConf.AddConfigPath("./conf/") //设置配置文件的搜索目录
	ViperConf.SetConfigType("toml")
	if err := ViperConf.ReadInConfig(); err != nil { // 加载配置文件内容
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	go ReloadConfig(ViperConf)
}

// 热加载
func ReloadConfig(ViperCon *viper.Viper) {
	ViperCon.WatchConfig()
	ViperCon.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("Detect config change: %s \n", in.String())
	})
}
