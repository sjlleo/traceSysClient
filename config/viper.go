package config

import "github.com/spf13/viper"

func InitConfig() {
	viper.AddConfigPath("./")          // 设置读取路径：就是在此路径下搜索配置文件。
	viper.AddConfigPath("./config/")   // 设置读取路径：就是在此路径下搜索配置文件。
	viper.SetConfigFile("config.yaml") // 设置被读取文件的全名，包括扩展名。
	// viper.SetDefault("backCallURL", "")
	// viper.SetDefault("token", "LeoCapstoneProject")
	viper.ReadInConfig() // 读取配置文件： 这一步将配置文件变成了 Go语言的配置文件对象包含了 map，string 等对象。
	viper.WriteConfig()
}
