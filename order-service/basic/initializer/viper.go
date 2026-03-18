package initializer

import (
	"fmt"
	"github.com/spf13/viper"
	"order/order-service/basic/config"
	"path/filepath"
)

func ViperInit() {
	var err error
	configPath := filepath.Join(GetProjectRoot(), "config.yaml")
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config.GlobalConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("配置加载成功")
}
