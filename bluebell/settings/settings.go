package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")          // 当前工作目录
	viper.AddConfigPath("./settings") // 你的 config.yaml 如果在 web_app 下面
	err = viper.ReadInConfig()        // 读取配置信息
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed!, err:%v\n", err) // 读取配置信息失败
		return
	}

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return
}
