package app

import (
	"fmt"
	"github.com/spf13/viper"
	"go-code/awesomeProject1/app/model"
	"go-code/awesomeProject1/app/router"
)

func Start() {
	model.NewMysql()
	model.NewRdb()
	defer func() {
		model.Close()
	}()
	port := 8080
	router.New(port)
}
func init() {
	// 设置配置文件名称(不含扩展名)
	viper.SetConfigName("config")
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 设置配置文件路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/app/")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误
			fmt.Println("No config file found")
		} else {
			// 其他错误
			fmt.Println("Error reading config file:", err)
		}
	}
}
