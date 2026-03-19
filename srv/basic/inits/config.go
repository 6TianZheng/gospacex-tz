package inits

import (
	"fmt"
	"github.com/spf13/viper"
	"gospacex-tz/srv/basic/config"
)

func ViperInit() {
	viper.SetConfigFile("../../../config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic("配置读取失败: " + err.Error())
	}

	err = viper.Unmarshal(&config.GlobalConfig)
	if err != nil {
		panic("配置解析失败: " + err.Error())
	}

	fmt.Println(config.GlobalConfig)
}
