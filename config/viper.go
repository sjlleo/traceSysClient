package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() {
	// viper.AddConfigPath("./")                      // 设置读取路径：就是在此路径下搜索配置文件。
	// viper.AddConfigPath("./config/")               // 设置读取路径：就是在此路径下搜索配置文件。
	// 配置文件名， 不加扩展
	viper.SetConfigName("config") // name of config file (without extension)
	// 设置文件的扩展名
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	// 查找配置文件所在路径
	viper.AddConfigPath("/usr/local/traceClient/") // path to look for the config file in
	// 在当前路径进行查找
	viper.AddConfigPath(".")         // optionally look for config in the working directory
	viper.AddConfigPath("./config/") // optionally look for config in the working directory

	// 开始查找并读取配置文件
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.ReadInConfig() // 读取配置文件： 这一步将配置文件变成了 Go语言的配置文件对象包含了 map，string 等对象。
	// viper.WriteConfig()
}
