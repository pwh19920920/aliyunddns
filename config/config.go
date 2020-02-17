package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var SystemConfig Config

type AliYunConfig struct {
	AliKeyId     string
	AliKeySecret string
	RecordDomain string
	RecordRr     string
	RecordType   string
	RecordId     string
	RecordTtl    int
}

type CronConfig struct {
	CronExp string
}

type LoggerConfig struct {
	TimestampFormat string
	LogPath         string
	LogLevel        string
}

type Config struct {
	AliYunConf AliYunConfig
	LoggerConf LoggerConfig
	CronConf   CronConfig
}

func LoadConfig() Config {
	//读取yaml文件
	v := viper.New()

	//设置读取的配置文件
	v.SetConfigName("config")

	//添加读取的配置文件路径
	v.AddConfigPath(".")

	//设置配置文件类型
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	return parseYaml(v)
}

func parseYaml(v *viper.Viper) Config {
	if err := v.Unmarshal(&SystemConfig); err != nil {
		fmt.Printf("err:%s", err)
	}

	return SystemConfig
}
