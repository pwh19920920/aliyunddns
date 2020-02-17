package main

import (
	"aliyunddns/aliyun"
	"aliyunddns/config"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 读取配置
	config.LoadConfig()

	// 初始化日志
	config.InitLogger()

	// 新建一个定时任务对象
	c := cron.New()
	c.AddFunc(config.SystemConfig.CronConf.CronExp, func() {
		// 更新阿里云ip
		aliyun.ExecuteUpdateAliYunDnsIp()
	})
	c.Start()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("服务优雅关闭")
}
